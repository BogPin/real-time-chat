package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BogPin/real-time-chat/backend/api/controllers"
	"github.com/BogPin/real-time-chat/backend/api/models/chat"
	"github.com/BogPin/real-time-chat/backend/api/models/message"
	"github.com/BogPin/real-time-chat/backend/api/models/participant"
	"github.com/BogPin/real-time-chat/backend/api/models/user"
	"github.com/BogPin/real-time-chat/backend/api/services"
	"github.com/BogPin/real-time-chat/backend/api/wss"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slices"
)

func main() {
	godotenv.Load()
	dbUser := getEnvVar("DB_USER")
	dbPassword := getEnvVar("DB_PASS")
	dbName := getEnvVar("DB_NAME")
	dbHost := getEnvVar("DB_HOST")
	dbPort := getEnvVar("DB_PORT")
	db, err := dbInit(dbUser, dbPassword, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authService := getEnvVar("AUTH_SERVICE")

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(controllers.GetAuthMiddleware(authService, controllers.GetTokenFromHeader))

	userStorer := user.UserStorer{DB: db}
	userService := services.UserService{UserStorer: &userStorer}
	usersRouter := apiRouter.PathPrefix("/users").Subrouter()
	controllers.RegisterUsersRoutes(usersRouter, userService)

	chatStorer := chat.ChatStorer{DB: db}
	chatService := services.ChatService{ChatStorer: &chatStorer}
	chatsRouter := apiRouter.PathPrefix("/chats").Subrouter()
	controllers.RegisterChatsRoutes(chatsRouter, chatService)

	messageStorer := message.MessageStorer{DB: db}
	messageService := services.MessageService{MessageStorer: &messageStorer}
	messagesRouter := apiRouter.PathPrefix("/messages").Subrouter()
	controllers.RegisterMessagesRoutes(messagesRouter, messageService)

	participantStorer := participant.ParticipantStorer{DB: db}
	participantService := services.ParticipantService{ParticipantStorer: &participantStorer}
	participantRouter := apiRouter.PathPrefix("/participants").Subrouter()
	controllers.RegisterParticipantRoutes(participantRouter, participantService)

	wsServer := wss.NewWsServer()
	wsRouter := router.PathPrefix("/ws").Subrouter()
	authMiddleware := controllers.GetAuthMiddleware(authService, controllers.GetTokenFromQuery)
	wsRouter.Use(authMiddleware)
	wsRouter.Path("").HandlerFunc(wsServer.HttpHandler).Methods("GET")

	wsServer.HandleConnection(func(socket *wss.Socket) {
		chats, err := chatService.GetUserChats(socket.UserId)
		if err != nil {
			socket.Disconnect(websocket.CloseInternalServerErr, err.Message())
		}
		for _, chat := range chats {
			socket.Join(chat.Id)
		}

		socket.On("message", func(data any) {
			msg := message.MessageFromRequest{}
			if err := mapstructure.Decode(data, &msg); err != nil {
				socket.Message(wss.NewErrorInvalidDataFormatMessage(msg))
			}
			i := slices.IndexFunc(chats, func(c chat.Chat) bool { return msg.ChatId == c.Id })
			if i == -1 {
				socket.Message(wss.NewErrorMessage("not allowed to write to that chat"))
				return
			}
			chatRoom, err := wsServer.Rooms.Get(msg.ChatId)
			if err != nil {
				socket.Message(wss.NewErrorMessage(err.Error()))
				return
			}
			fullMessage, httpErr := messageService.Create(socket.UserId, msg)
			if httpErr != nil {
				socket.Message(wss.NewErrorMessage(httpErr.Message()))
				return
			}
			chatRoom.Send(socket.UserId, wss.NewMessage("message", fullMessage))
			socket.Message(wss.NewMessage("message", fullMessage))
		})

		socket.On("disconnect", func(data any) {
			chatRooms := wsServer.Rooms.GetAllForUser(socket.UserId)
			for _, chatRoom := range chatRooms {
				msg := message.MessageDTO{
					SenderId: socket.UserId,
					ChatId:   chatRoom.Id,
					Type:     "text",
					Content:  "bye, I leave",
				}
				chatRoom.Send(socket.UserId, wss.NewMessage("message", msg))
			}
		})

	})

	port := ":" + getEnvVar("PORT")
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}

func dbInit(user, password, host, port, dbname string) (*sql.DB, error) {
	conStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname)
	return sql.Open("postgres", conStr)
}

func getEnvVar(name string) string {
	variable, present := os.LookupEnv(name)
	if !present {
		log.Fatalf("%s env variable is missing", name)
	}
	return variable
}
