package wss

import (
	"log"

	"github.com/gorilla/websocket"
)

type Socket struct {
	UserId    int
	conn      *websocket.Conn
	listeners map[string][]func(data any)
	server    *WsServer
}

func NewSocket(userId int, conn *websocket.Conn, server *WsServer) *Socket {
	return &Socket{
		UserId:    userId,
		conn:      conn,
		listeners: make(map[string][]func(data any)),
		server:    server,
	}
}

func (s *Socket) PostDisconnect() {
	s.server.Conns.Remove(s.UserId)
	chats := s.server.Rooms.GetAllForUser(s.UserId)
	for _, chat := range chats {
		chat.Remove(s.UserId)
	}
}

func (s *Socket) Disconnect(code int, reason string) {
	cm := websocket.FormatCloseMessage(code, reason)
	err := s.conn.WriteMessage(websocket.CloseMessage, cm)
	if err != nil {
		log.Println(err)
	}
	s.conn.Close()
	s.PostDisconnect()
}

func (s *Socket) On(event string, listener func(data any)) {
	s.listeners[event] = append(s.listeners[event], listener)
}

func (s *Socket) emit(event string, data any) {
	if listeners, ok := s.listeners[event]; ok {
		for _, listener := range listeners {
			listener(data)
		}
	}
}

func (s *Socket) Message(msg Message) error {
	return s.conn.WriteJSON(msg)
}

func (s *Socket) Join(roomId int) {
	room, err := s.server.Rooms.Get(roomId)
	if err != nil {
		room = s.server.Rooms.Add(roomId)
	}
	room.Add(s.UserId, s.conn)
}

func (s *Socket) Leave(roomId int) error {
	room, err := s.server.Rooms.Get(roomId)
	if err != nil {
		return err
	}
	room.Remove(s.UserId)
	return nil
}
