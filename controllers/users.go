package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BogPin/real-time-chat/models/user"
	"github.com/BogPin/real-time-chat/services"
	"github.com/gorilla/mux"
)

func RegisterUsersRoutes(router *mux.Router, service services.User) {
	router.Path("/").HandlerFunc(getUser(service)).Methods("GET")
	router.Path("/").HandlerFunc(createUser(service)).Methods("POST")
	router.Path("/").HandlerFunc(updateUser(service)).Methods("PATCH")
	router.Path("/").HandlerFunc(deleteUser(service)).Methods("DELETE")
	router.Path("/{id}").HandlerFunc(getUser(service)).Methods("GET")
	router.Path("/{id}").HandlerFunc(deleteUser(service)).Methods("DELETE")
}

func getUser(service services.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlId, ok := mux.Vars(r)["id"]
		queryId := r.URL.Query().Get("id")
		if ok && queryId != "" && urlId != queryId {
			http.Error(w, "id in url and id in query are incompatible", http.StatusBadRequest)
			return
		}
		var id int
		var err error
		if ok {
			id, err = strconv.Atoi(urlId)
		} else {
			id, err = strconv.Atoi(queryId)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := service.GetOne(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func createUser(service services.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userDto user.UserDTO
		err := json.NewDecoder(r.Body).Decode(&userDto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := service.Create(userDto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func updateUser(service services.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user user.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newUser, err := service.Update(user)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func deleteUser(service services.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlId, ok := mux.Vars(r)["id"]
		queryId := r.URL.Query().Get("id")
		if ok && queryId != "" && urlId != queryId {
			http.Error(w, "id in url and id in query are incompatible", http.StatusBadRequest)
			return
		}
		var id int
		var err error
		if ok {
			id, err = strconv.Atoi(urlId)
		} else {
			id, err = strconv.Atoi(queryId)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := service.Delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
