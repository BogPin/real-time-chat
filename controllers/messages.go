package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BogPin/real-time-chat/models/message"
	"github.com/BogPin/real-time-chat/services"
	"github.com/gorilla/mux"
)

func RegisterMessagesRoutes(router *mux.Router, service services.Message) {
	router.Path("/").HandlerFunc(createMessage(service)).Methods("POST")
	router.Path("/").HandlerFunc(updateMessage(service)).Methods("PATCH")
	router.Path("/{id}").HandlerFunc(getMessage(service)).Methods("GET")
	router.Path("/{id}").HandlerFunc(deleteMessage(service)).Methods("DELETE")
}

func getMessage(service services.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		message, err := service.GetMessage(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func deleteMessage(service services.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		message, err := service.Delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func updateMessage(service services.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var message message.Message
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newMessage, err := service.Update(message)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(newMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func createMessage(service services.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var messageDTO message.MessageDTO
		err := json.NewDecoder(r.Body).Decode(&messageDTO)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		message, err := service.Create(messageDTO)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
