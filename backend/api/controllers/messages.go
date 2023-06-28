package controllers

import (
	"net/http"
	"strconv"

	"github.com/BogPin/real-time-chat/backend/api/services"
	"github.com/BogPin/real-time-chat/backend/api/utils"
	"github.com/gorilla/mux"
)

func RegisterMessagesRoutes(router *mux.Router, service services.IMessageService) {
	router.Path("/{id}").HandlerFunc(getMessage(service)).Methods("GET")
	router.Path("").HandlerFunc(getMessages(service)).Methods("GET")
}

func getMessage(service services.IMessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messageId, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			WriteError(w, utils.NewHttpError(err, http.StatusBadRequest))
			return
		}

		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		message, httpErr := service.GetOne(payload.UserId, messageId)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, message)
	}
}

func getMessages(service services.IMessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId, err := strconv.Atoi(r.URL.Query().Get("chatId"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		messages, httpErr := service.GetChatMessages(payload.UserId, chatId, page)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, messages)
	}
}
