package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/auth/models/user"
	"github.com/BogPin/real-time-chat/backend/auth/services"
	"github.com/BogPin/real-time-chat/backend/auth/utils"
	"github.com/gorilla/mux"
)

type LoginEndpoint struct {
	method   string
	name     string
	service  services.UserService
	jwtStrat utils.JWTStrategy
}

func NewLoginEndpoint(method string, name string, service services.UserService, jwtStrat utils.JWTStrategy) LoginEndpoint {
	return LoginEndpoint{method, name, service, jwtStrat}
}

func (le LoginEndpoint) Add(router *mux.Router) {
	router.Path(le.name).HandlerFunc(le.Handle).Methods(le.method)
}

func (le LoginEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	var creds user.Credentials
	jsonErr := json.NewDecoder(r.Body).Decode(&creds)
	if jsonErr != nil {
		writeError(w, utils.NewHttpError(jsonErr, http.StatusBadRequest))
		return
	}
	user, err := le.service.ValidateUser(creds)
	if err != nil {
		writeError(w, err)
		return
	}
	token, err := le.jwtStrat.CreateJWT(user.Id)
	if err != nil {
		writeError(w, err)
		return
	}
	resp := tokenBody{token}
	jsonErr = json.NewEncoder(w).Encode(resp)
	if jsonErr != nil {
		writeError(w, utils.NewHttpError(jsonErr, http.StatusInternalServerError))
		return
	}
}
