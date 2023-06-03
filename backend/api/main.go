package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BogPin/real-time-chat/backend/api/models/chat"

	"github.com/BogPin/real-time-chat/backend/api/controllers"
	"github.com/BogPin/real-time-chat/backend/api/models/message"
	"github.com/BogPin/real-time-chat/backend/api/models/user"
	"github.com/BogPin/real-time-chat/backend/api/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	db, err := dbInit(dbUser, dbPassword, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	userStorer := user.UserStorer{DB: db}
	userService := services.UserService{UserStorer: &userStorer}
	usersRouter := router.PathPrefix("/users").Subrouter()
	controllers.RegisterUsersRoutes(usersRouter, userService)

	chatStorer := chat.ChatStorer{DB: db}
	chatService := services.ChatService{ChatStorer: &chatStorer}
	chatsRouter := router.PathPrefix("/chats").Subrouter()
	controllers.RegisterChatsRoutes(chatsRouter, chatService)

	messageStorer := message.MessageStorer{DB: db}
	messageService := services.MessageService{MessageStorer: &messageStorer}
	messagesRouter := router.PathPrefix("/messages").Subrouter()
	controllers.RegisterMessagesRoutes(messagesRouter, messageService)
	http.ListenAndServe(":8080", router)
}

func dbInit(user, password, host, port, dbname string) (*sql.DB, error) {
	conStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname)
	return sql.Open("postgres", conStr)
}
