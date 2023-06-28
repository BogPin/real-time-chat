package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BogPin/real-time-chat/backend/api/models"
	"github.com/BogPin/real-time-chat/backend/api/services"
	"github.com/BogPin/real-time-chat/backend/api/utils"
	"github.com/gorilla/mux"
)

func RegisterChatsRoutes(router *mux.Router, service services.IChatService) {
	router.Path("").HandlerFunc(createChat(service)).Methods("POST")
	router.Path("").HandlerFunc(getChats(service)).Methods("GET")
	router.Path("/{id}").HandlerFunc(getChat(service)).Methods("GET")
	router.Path("/{id}").HandlerFunc(updateChat(service)).Methods("PATCH")
	router.Path("/{id}").HandlerFunc(deleteChat(service)).Methods("DELETE")
}

func createChat(service services.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var chatWithTitle models.ChatFromRequest
		err := json.NewDecoder(r.Body).Decode(&chatWithTitle)
		if err != nil {
			WriteError(w, utils.NewHttpError(err, http.StatusBadRequest))
			return
		}

		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		chat, httpErr := service.Create(payload.UserId, chatWithTitle)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, chat)
	}
}

func getChats(service services.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		chats, httpErr := service.GetUserChats(payload.UserId)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, chats)
	}
}

func getChat(service services.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			WriteError(w, utils.NewHttpError(err, http.StatusBadRequest))
			return
		}

		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		chat, httpErr := service.GetOne(payload.UserId, chatId)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, chat)
	}
}

func updateChat(service services.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var chat models.Chat
		err := json.NewDecoder(r.Body).Decode(&chat)
		if err != nil {
			WriteError(w, utils.NewHttpError(err, http.StatusBadRequest))
			return
		}

		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		updChat, httpErr := service.Update(payload.UserId, chat)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, updChat)
	}
}

func deleteChat(service services.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			WriteError(w, utils.NewHttpError(err, http.StatusBadRequest))
			return
		}

		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		dltChat, httpErr := service.Delete(payload.UserId, chatId)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, dltChat)
	}
}
