package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BogPin/real-time-chat/backend/auth/controllers"
	"github.com/BogPin/real-time-chat/backend/auth/models/user"
	"github.com/BogPin/real-time-chat/backend/auth/services"
	"github.com/BogPin/real-time-chat/backend/auth/utils"
	"github.com/golang-jwt/jwt"
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET env variable is missing")
	}
	jwtStrat := utils.NewJWTStrategy(jwtSecret, jwt.GetSigningMethod("HS256"))

	router := mux.NewRouter().PathPrefix("/auth").Subrouter()

	userStorer := user.Storer{DB: db}
	userService := services.UserService{UserStorer: userStorer}

	controllers.NewLoginEndpoint("POST", "/login", userService, jwtStrat).Add(router)
	controllers.NewRegisterEndpoint("POST", "/register", userService, jwtStrat).Add(router)
	controllers.NewValidateTokenEndpoint("POST", "/validate", jwtStrat).Add(router)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env variable is missing")
	}

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}

func dbInit(user, password, host, port, dbname string) (*sql.DB, error) {
	conStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname)
	return sql.Open("postgres", conStr)
}
