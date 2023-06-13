package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BogPin/real-time-chat/backend/api/models/participant"
	"github.com/BogPin/real-time-chat/backend/api/services"
	"github.com/BogPin/real-time-chat/backend/api/utils"
	"github.com/gorilla/mux"
)

func RegisterParticipantRoutes(router *mux.Router, service services.Participant) {
	router.Path("").HandlerFunc(createParticipant(service)).Methods("POST")
	router.Path("/{id}").HandlerFunc(updateParticipant(service)).Methods("PATCH")
	router.Path("").HandlerFunc(deleteParticipant(service)).Methods("DELETE")
}

func createParticipant(service services.Participant) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		part := participant.Participant{
			Role: "member",
		}
		err := json.NewDecoder(r.Body).Decode(&part)
		if err != nil {
			WriteError(w, utils.NewHttpError(err, http.StatusBadRequest))
			return
		}

		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		participant, httpErr := service.Create(payload.UserId, part)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, participant)
	}
}

func updateParticipant(service services.Participant) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var part participant.Participant
		err := json.NewDecoder(r.Body).Decode(&part)
		if err != nil {
			WriteError(w, utils.NewHttpError(err, http.StatusBadRequest))
			return
		}

		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		updPart, httpErr := service.Update(payload.UserId, part)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, updPart)
	}
}

func deleteParticipant(service services.Participant) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		chatId, err := strconv.Atoi(r.URL.Query().Get("chatId"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		part := participant.Participant{
			UserId: userId,
			ChatId: chatId,
		}

		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		dltPart, httpErr := service.Delete(payload.UserId, part)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, dltPart)
	}
}
