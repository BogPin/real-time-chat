package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/api/models"
	"github.com/BogPin/real-time-chat/backend/api/services"
	"github.com/BogPin/real-time-chat/backend/api/utils"
	"github.com/gorilla/mux"
)

func RegisterUsersRoutes(router *mux.Router, service services.IUserService) {
	router.Path("").HandlerFunc(getUser(service)).Methods("GET")
	router.Path("").HandlerFunc(updateUser(service)).Methods("PATCH")
	router.Path("").HandlerFunc(deleteUser(service)).Methods("DELETE")
}

func getUser(service services.IUserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		user, httpErr := service.GetOne(payload.UserId)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, user)
	}
}

func updateUser(service services.IUserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			WriteError(w, utils.NewHttpError(err, http.StatusBadRequest))
			return
		}

		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		newUser, httpErr := service.Update(payload.UserId, user)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, newUser)
	}
}

func deleteUser(service services.IUserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			WriteError(w, utils.NewHttpError(err, http.StatusBadRequest))
			return
		}

		payload, ok := r.Context().Value(TokenPayloadKey).(TokenPayload)
		if !ok {
			WriteError(w, ErrNoUserPayloadInContext)
			return
		}

		dltUser, httpErr := service.Delete(payload.UserId, user)
		if httpErr != nil {
			WriteError(w, httpErr)
			return
		}

		writeResponce(w, dltUser)
	}
}
