package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/auth/models/user"
	"github.com/BogPin/real-time-chat/backend/auth/utils"
)

type UserService struct {
	UserStorer user.Storer
}

func (us UserService) ValidateUser(creds user.Credentials) (*user.User, utils.HttpError) {
	usr, err := us.UserStorer.GetUser(creds.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := fmt.Errorf("no user with name %s", creds.Username)
			return nil, utils.NewHttpError(err, http.StatusNotFound)
		}
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}
	if usr.PasswordHash != creds.PasswordHash {
		return nil, utils.NewHttpError(errors.New("incorrect password"), http.StatusUnauthorized)
	}
	return &usr, nil
}

func (us UserService) RegisterUser(creds user.Credentials) (*user.User, utils.HttpError) {
	_, err := us.UserStorer.GetUser(creds.Username)
	if err == nil {
		msg := fmt.Sprintf("username %s is already taken", creds.Username)
		return nil, utils.NewHttpError(errors.New(msg), http.StatusBadRequest)
	}
	if err != sql.ErrNoRows {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}
	usr, err := us.UserStorer.CreateUser(creds)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}
	return &usr, nil
}
