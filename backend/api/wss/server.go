package wss

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"

	"github.com/BogPin/real-time-chat/backend/api/controllers"
	"github.com/BogPin/real-time-chat/backend/api/utils"
	"github.com/gorilla/websocket"
)

const MessageFormatErr = "{\"error\":\"message must be in format: { event: string, data: any }\"}"

type Message struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

func NewMessage(event string, data any) Message {
	return Message{event, data}
}

func NewErrorMessage(errMsg string) Message {
	return Message{Event: "error", Data: errMsg}
}

func NewErrorInvalidDataFormatMessage(val any) Message {
	t := reflect.TypeOf(val)
	dataFormat := utils.GetJSONSignature(t)
	errMsg := fmt.Sprintf("message data must be in format: %s", dataFormat)
	return Message{Event: "error", Data: errMsg}
}

type safeConns struct {
	mu    sync.RWMutex
	conns map[int]*websocket.Conn
}

func (sc *safeConns) Get(userId int) (*websocket.Conn, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	if conn, ok := sc.conns[userId]; ok {
		return conn, nil
	}
	return nil, fmt.Errorf("no socket with user id %d", userId)
}

func (sc *safeConns) SendAll(fromUser int, msg Message) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	for userId, conn := range sc.conns {
		if userId != fromUser {
			if err := conn.WriteJSON(msg); err != nil {
				fmt.Printf("error while sending all msg %v from user %v:\n", msg, fromUser)
			}
		}
	}
}

func (sc *safeConns) Add(userId int, conn *websocket.Conn) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.conns[userId] = conn
}

func (sc *safeConns) Remove(userId int) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	delete(sc.conns, userId)
}

type WsServer struct {
	Conns         safeConns
	Rooms         safeRooms
	upgrader      websocket.Upgrader
	socketHandler func(socket *Socket)
}

func NewWsServer() *WsServer {
	return &WsServer{
		Conns: safeConns{
			conns: make(map[int]*websocket.Conn),
		},
		Rooms: safeRooms{
			rooms: make([]*room, 0),
		},
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (wss *WsServer) HandleConnection(handler func(socket *Socket)) {
	wss.socketHandler = handler
}

func (wss *WsServer) listenMessages(socket *Socket) {
	for {
		messageType, msg, err := socket.conn.ReadMessage()
		if err != nil {
			fmt.Println("read message error:", err)
			socket.emit("disconnect", nil)
			socket.PostDisconnect()
			return
		}
		if messageType != websocket.TextMessage {
			data := []byte("only text messages are allowed")
			socket.conn.WriteMessage(websocket.TextMessage, data)
			continue
		}
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			socket.conn.WriteMessage(websocket.TextMessage, []byte(MessageFormatErr))
		}

		go socket.emit(message.Event, message.Data)
	}
}

func (wss *WsServer) HttpHandler(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value(controllers.TokenPayloadKey).(controllers.TokenPayload)
	if !ok {
		controllers.WriteError(w, controllers.ErrNoUserPayloadInContext)
		return
	}

	conn, _ := wss.upgrader.Upgrade(w, r, nil)
	socket := NewSocket(payload.UserId, conn, wss)
	wss.Conns.Add(socket.UserId, conn)
	wss.socketHandler(socket)
	go wss.listenMessages(socket)
}
