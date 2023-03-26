package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/BogPin/real-time-chat/entities"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi user")
}

func main() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	host := "localhost"
	db, err := dbInit(user, password, host, dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", handler)
	http.HandleFunc("/user", handleUser)
	http.ListenAndServe(":8080", nil)

}

func dbInit(user, password, host, dbname string) (*sql.DB, error) {
	conStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", user, password, host, dbname)
	return sql.Open("postgres", conStr)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Обработка GET запроса...

	case http.MethodPost:
		userDto := entities.UserDTO{}
		err := json.NewDecoder(r.Body).Decode(&userDto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := entities.

	case http.MethodPatch:
		//

	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
