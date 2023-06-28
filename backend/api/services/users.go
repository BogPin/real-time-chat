package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/api/models"
	"github.com/BogPin/real-time-chat/backend/api/utils"
)

type IUserService interface {
	GetOne(userId int) (*models.User, utils.HttpError)
	Update(userId int, user models.User) (*models.User, utils.HttpError)
	Delete(userId int, user models.User) (*models.User, utils.HttpError)
}

type UserService struct {
	UserStorer models.IUserStorer
}

func NewUserService(userStorer models.IUserStorer) UserService {
	return UserService{UserStorer: userStorer}
}

func (us UserService) GetOne(userId int) (*models.User, utils.HttpError) {
	user, err := us.UserStorer.GetOne(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := fmt.Errorf("no user with id %d", userId)
			return nil, utils.NewHttpError(err, http.StatusNotFound)
		}
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}
	return user, nil
}

func (us UserService) Update(userId int, user models.User) (*models.User, utils.HttpError) {
	//TODO: update only non nullish fields
	if userId != user.Id {
		err := errors.New("cannot update other users")
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	newUser, err := us.UserStorer.Update(user)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	return newUser, nil
}

func (us UserService) Delete(userId int, user models.User) (*models.User, utils.HttpError) {
	if userId != user.Id {
		err := errors.New("cannot delete other users")
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	deletedUser, err := us.UserStorer.Delete(userId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	return deletedUser, nil
}
