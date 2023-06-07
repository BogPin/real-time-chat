package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/BogPin/real-time-chat/backend/api/models/user"
	"github.com/BogPin/real-time-chat/backend/api/services"
	"github.com/gorilla/mux"
)

func RegisterUsersRoutes(router *mux.Router, service services.User) {
	router.Path("/").HandlerFunc(updateUser(service)).Methods("PATCH")
	router.Path("/{id}").HandlerFunc(getUser(service)).Methods("GET")
	router.Path("/{id}").HandlerFunc(deleteUser(service)).Methods("DELETE")
}

func getUser(service services.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
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
			if errors.Is(err, sql.ErrNoRows) {
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
		id, err := strconv.Atoi(mux.Vars(r)["id"])
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
