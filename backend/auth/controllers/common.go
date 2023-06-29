package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/auth/utils"
	"github.com/gorilla/mux"
)

type TokenBody struct {
	Token string `json:"token"`
}

type errorResponce struct {
	ErrorMsg string `json:"error"`
}

type Endpoint interface {
	Handle(http.ResponseWriter, *http.Request)
	Add(*mux.Router)
}

func WriteError(w http.ResponseWriter, err utils.HttpError) {
	w.WriteHeader(err.Status())
	json.NewEncoder(w).Encode(errorResponce{err.Message()})
}
