package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/api/controllers"
	"github.com/BogPin/real-time-chat/backend/api/models"
	"github.com/BogPin/real-time-chat/backend/api/services"
	"github.com/BogPin/real-time-chat/backend/api/utils"
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
	dbUser := utils.GetEnvVar("DB_USER")
	dbPassword := utils.GetEnvVar("DB_PASS")
	dbName := utils.GetEnvVar("DB_NAME")
	dbHost := utils.GetEnvVar("DB_HOST")
	dbPort := utils.GetEnvVar("DB_PORT")
	db, err := dbInit(dbUser, dbPassword, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authService := utils.GetEnvVar("AUTH_SERVICE")

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(controllers.GetAuthMiddleware(authService, controllers.GetTokenFromHeader))

	userStorer := models.NewUserStorer(db)
	userService := services.NewUserService(userStorer)
	usersRouter := apiRouter.PathPrefix("/users").Subrouter()
	controllers.RegisterUsersRoutes(usersRouter, userService)

	participantStorer := models.NewParticipantStorer(db)
	participantService := services.NewParticipantService(participantStorer)
	participantRouter := apiRouter.PathPrefix("/participants").Subrouter()
	controllers.RegisterParticipantRoutes(participantRouter, participantService)

	messageStorer := models.NewMessageStorer(db)
	messageService := services.NewMessageService(messageStorer, participantService)
	messagesRouter := apiRouter.PathPrefix("/messages").Subrouter()
	controllers.RegisterMessagesRoutes(messagesRouter, messageService)

	chatStorer := models.NewChatStorer(db)
	chatService := services.NewChatService(chatStorer, participantStorer, messageStorer, participantService)
	chatsRouter := apiRouter.PathPrefix("/chats").Subrouter()
	controllers.RegisterChatsRoutes(chatsRouter, chatService)

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
			msg := models.MessageFromRequest{}
			if err := mapstructure.Decode(data, &msg); err != nil {
				socket.Message(wss.NewErrorInvalidDataFormatMessage(msg))
			}
			i := slices.IndexFunc(chats, func(c models.Chat) bool { return msg.ChatId == c.Id })
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
				msg := models.MessageDTO{
					SenderId: socket.UserId,
					ChatId:   chatRoom.Id,
					Type:     "text",
					Content:  "bye, I leave",
				}
				chatRoom.Send(socket.UserId, wss.NewMessage("message", msg))
			}
		})

	})

	port := ":" + utils.GetEnvVar("PORT")
	log.Printf("listening on %s", port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}

func dbInit(user, password, host, port, dbname string) (*sql.DB, error) {
	conStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname)
	return sql.Open("postgres", conStr)
}
