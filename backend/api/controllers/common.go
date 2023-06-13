package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/api/utils"
	"github.com/gorilla/mux"
)

type TokenPayload struct {
	UserId int `json:"userId"`
}

type errorResponce struct {
	ErrorMsg string `json:"error"`
}

type Endpoint interface {
	Handle(http.ResponseWriter, *http.Request)
	Add(*mux.Router)
}

var ErrNoUserPayloadInContext = utils.NewHttpError(
	errors.New("could not get user payload from context"),
	http.StatusInternalServerError,
)

func WriteError(w http.ResponseWriter, err utils.HttpError) {
	w.WriteHeader(err.Status())
	json.NewEncoder(w).Encode(errorResponce{err.Message()})
}

func writeResponce(w http.ResponseWriter, resp any) {
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		WriteError(w, utils.NewHttpError(err, http.StatusInternalServerError))
		return
	}
}
