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

	jwtSecret := getEnvVar("JWT_SECRET")
	jwtStrat := utils.NewJWTStrategy(jwtSecret, jwt.GetSigningMethod("HS256"))

	router := mux.NewRouter().PathPrefix("/auth").Subrouter()

	userStorer := user.Storer{DB: db}
	userService := services.UserService{UserStorer: userStorer}

	controllers.NewLoginEndpoint("POST", "/login", userService, jwtStrat).Add(router)
	controllers.NewRegisterEndpoint("POST", "/register", userService, jwtStrat).Add(router)
	controllers.NewValidateTokenEndpoint("POST", "/validate", jwtStrat).Add(router)

	port := getEnvVar("PORT")
	log.Printf("listening on %s", port)
	err = http.ListenAndServe(":"+port, router)
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
