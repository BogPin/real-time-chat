package controllers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	apiControllers "github.com/BogPin/real-time-chat/backend/api/controllers"
	authControllers "github.com/BogPin/real-time-chat/backend/auth/controllers"
	authUser "github.com/BogPin/real-time-chat/backend/auth/models/user"
	authServices "github.com/BogPin/real-time-chat/backend/auth/services"
	authUtils "github.com/BogPin/real-time-chat/backend/auth/utils"
	"github.com/golang-jwt/jwt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	DB_PASSWORD  = "admin"
	JWT_SECRET   = "secret"
	EXPOSED_PORT = "5432"
	HOST_PORT    = "5433"
)

var db *sql.DB
var apiServer *httptest.Server
var authServer *httptest.Server
var authTokenN, authTokenB string

type TokenBody struct {
	Token string `json:"token"`
}

func getToken(cred authUser.Credentials) (string, error) {
	jsonBody, _ := json.Marshal(cred)
	loginEndpoint := fmt.Sprintf("%s/auth/login", authServer.URL)
	respToken, err := http.Post(loginEndpoint, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("Failed to create token: %v", err)
	}
	var tokenBody TokenBody
	err = json.NewDecoder(respToken.Body).Decode(&tokenBody)
	if err != nil {
		return "", fmt.Errorf("Failed to decode token response: %v", err)
	}
	return tokenBody.Token, nil
}

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "postgres",
		Tag:          "14.4",
		Env:          []string{fmt.Sprintf("POSTGRES_PASSWORD=%s", DB_PASSWORD)},
		ExposedPorts: []string{EXPOSED_PORT},
		PortBindings: map[docker.Port][]docker.PortBinding{
			EXPOSED_PORT: {
				{HostIP: "0.0.0.0", HostPort: HOST_PORT},
			},
		},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Tell docker to hard kill the container in 30 seconds
	if err := resource.Expire(30); err != nil {
		log.Fatalf("Could not set autokill of container: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	retryFunc := func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("port=%s user=postgres password=%s sslmode=disable", HOST_PORT, DB_PASSWORD))
		if err != nil {
			return err
		}
		return db.Ping()
	}
	if err := pool.Retry(retryFunc); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	mig, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = mig.Up(); err != nil {
		log.Fatalf("Error while executing migration: %s", err)
	}

	// Setup auth server
	authServerRouter := mux.NewRouter()
	authRouter := authServerRouter.PathPrefix("/auth").Subrouter()

	userStorer := authUser.Storer{DB: db}
	userService := authServices.UserService{UserStorer: userStorer}
	jwtStrat := authUtils.NewJWTStrategy(JWT_SECRET, jwt.GetSigningMethod("HS256"))

	authControllers.NewLoginEndpoint("POST", "/login", userService, jwtStrat).Add(authRouter)
	authControllers.NewRegisterEndpoint("POST", "/register", userService, jwtStrat).Add(authRouter)
	authControllers.NewValidateTokenEndpoint("POST", "/validate", jwtStrat).Add(authRouter)

	authServer = httptest.NewServer(authServerRouter)
	log.Println("auth server is listening on", authServer.URL)

	authTokenN, err = getToken(authUser.Credentials{Username: "holdennekt", PasswordHash: "abc"})
	if err != nil {
		log.Fatal(err)
	}

	authTokenB, err = getToken(authUser.Credentials{Username: "bogpin", PasswordHash: "abcd"})
	if err != nil {
		log.Fatal(err)
	}

	// Setup api server
	apiServerRouter := mux.NewRouter()
	apiRouter := apiServerRouter.PathPrefix("/api").Subrouter()
	urlSchema, err := url.Parse(authServer.URL)
	if err != nil {
		log.Fatalf("Error during parsing URL: %s", err)
	}
	apiRouter.Use(apiControllers.GetAuthMiddleware(urlSchema.Host, apiControllers.GetTokenFromHeader))
	chatsRouter := apiRouter.PathPrefix("/chats").Subrouter()
	chatsService := getChatService(db)
	apiControllers.RegisterChatsRoutes(chatsRouter, chatsService)

	apiServer = httptest.NewServer(apiServerRouter)
	log.Println("api server is listening on", apiServer.URL)

	// Run all tests
	code := m.Run()

	// Clean up after tests
	apiServer.Close()
	authServer.Close()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
