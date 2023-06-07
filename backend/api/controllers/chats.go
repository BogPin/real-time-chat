package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/BogPin/real-time-chat/backend/api/models/chat"
	"github.com/BogPin/real-time-chat/backend/api/services"
	"github.com/gorilla/mux"
)

func RegisterChatsRoutes(router *mux.Router, service services.Chat) {
	router.Path("/").HandlerFunc(createChat(service)).Methods("POST")
	router.Path("/").HandlerFunc(updateChat(service)).Methods("PATCH")
	router.Path("/{id}").HandlerFunc(getChat(service)).Methods("GET")
	router.Path("/{id}").HandlerFunc(deleteChat(service)).Methods("DELETE")
}

func getChat(service services.Chat) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		chat, err := service.GetChat(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(chat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func deleteChat(service services.Chat) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		chat, err := service.Delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(chat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func updateChat(service services.Chat) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var chat chat.Chat
		err := json.NewDecoder(r.Body).Decode(&chat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newChat, err := service.Update(chat)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(newChat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func createChat(service services.Chat) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var chatDTO chat.ChatDTO
		err := json.NewDecoder(r.Body).Decode(&chatDTO)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		chat, err := service.Create(chatDTO)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(chat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
