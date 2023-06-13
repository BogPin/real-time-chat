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

	router := mux.NewRouter()
	router.Use(controllers.GetAuthMiddleware(controllers.GetTokenFromHeader))

	apiRouter := router.PathPrefix("/api").Subrouter()

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

	port := ":" + getEnvVar("PORT")
	go func() {
		err = http.ListenAndServe(port, router)
		if err != nil {
			log.Fatal(err)
		}
	}()
	wsServer := wss.NewWsServer()
	wsServer.HandleConnection(func(socket *wss.Socket) {
		chats, err := chatService.GetUserChats(socket.UserId)
		if err != nil {
			socket.Disconnect(websocket.CloseInternalServerErr, err.Message())
		}
		for _, chat := range chats {
			socket.Join(chat.Id)
		}

		socket.On("message", func(data interface{}) {
			msg := message.MessageFromRequest{}
			if err := mapstructure.Decode(data, &msg); err != nil {
				socket.Message(wss.NewErrorInvalidDataFormatMessage(msg))
			}
			i := slices.IndexFunc(chats, func(c chat.Chat) bool { return msg.ChatId == c.Id })
			if i == -1 {
				socket.Message(wss.NewErrorMessage("not allowed to write to that chat"))
				return
			}
			chat, err := wsServer.Rooms.Get(msg.ChatId)
			if err != nil {
				socket.Message(wss.NewErrorMessage(err.Error()))
				return
			}
			fullMessage, httpErr := messageService.Create(socket.UserId, msg)
			if httpErr != nil {
				socket.Message(wss.NewErrorMessage(httpErr.Message()))
				return
			}
			chat.Send(socket.UserId, wss.NewMessage("message", fullMessage))
			socket.Message(wss.NewMessage("message", fullMessage))
		})

	})

	err = wsServer.ListenAndServe(":8082", "/ws")
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
