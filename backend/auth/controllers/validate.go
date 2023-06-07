package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/auth/utils"
	"github.com/gorilla/mux"
)

type ValidateTokenEndpoint struct {
	method   string
	name     string
	jwtStrat utils.JWTStrategy
}

func NewValidateTokenEndpoint(method string, name string, jwtStrat utils.JWTStrategy) ValidateTokenEndpoint {
	return ValidateTokenEndpoint{method, name, jwtStrat}
}

func (vte ValidateTokenEndpoint) Add(router *mux.Router) {
	router.Path(vte.name).HandlerFunc(vte.Handle).Methods(vte.method)
}

func (vte ValidateTokenEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	var body tokenBody
	jsonErr := json.NewDecoder(r.Body).Decode(&body)
	if jsonErr != nil {
		writeError(w, utils.NewHttpError(jsonErr, http.StatusBadRequest))
		return
	}
	tokenPayload, err := vte.jwtStrat.DecodeJWT(body.Token)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonErr = json.NewEncoder(w).Encode(tokenPayload)
	if jsonErr != nil {
		writeError(w, utils.NewHttpError(jsonErr, http.StatusInternalServerError))
		return
	}
}
