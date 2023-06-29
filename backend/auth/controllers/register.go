package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/auth/models/user"
	"github.com/BogPin/real-time-chat/backend/auth/services"
	"github.com/BogPin/real-time-chat/backend/auth/utils"
	"github.com/gorilla/mux"
)

type RegisterEndpoint struct {
	method   string
	name     string
	service  services.UserService
	jwtStrat utils.JWTStrategy
}

func NewRegisterEndpoint(method string, name string, service services.UserService, jwtStrat utils.JWTStrategy) RegisterEndpoint {
	return RegisterEndpoint{method, name, service, jwtStrat}
}

func (re RegisterEndpoint) Add(router *mux.Router) {
	router.Path(re.name).HandlerFunc(re.Handle).Methods(re.method)
}

func (re RegisterEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	var creds user.Credentials
	jsonErr := json.NewDecoder(r.Body).Decode(&creds)
	if jsonErr != nil {
		WriteError(w, utils.NewHttpError(jsonErr, http.StatusBadRequest))
		return
	}
	user, err := re.service.RegisterUser(creds)
	if err != nil {
		WriteError(w, err)
		return
	}
	token, err := re.jwtStrat.CreateJWT(user.Id)
	if err != nil {
		WriteError(w, err)
		return
	}
	resp := TokenBody{token}
	jsonErr = json.NewEncoder(w).Encode(resp)
	if jsonErr != nil {
		WriteError(w, utils.NewHttpError(jsonErr, http.StatusInternalServerError))
		return
	}
}
