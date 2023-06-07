package services

import (
	"github.com/BogPin/real-time-chat/backend/api/models/user"
)

type User interface {
	GetAll() ([]user.User, error)
	GetOne(id int) (*user.User, error)
	Update(user user.User) (*user.User, error)
	Delete(id int) (*user.User, error)
}

type UserService struct {
	UserStorer *user.UserStorer
}

func (us UserService) GetAll() ([]user.User, error) {
	return us.UserStorer.GetAll()
}

func (us UserService) GetOne(id int) (*user.User, error) {
	return us.UserStorer.GetOne(id)
}

func (us UserService) Update(user user.User) (*user.User, error) {
	//TODO: update only non nullish fields
	return us.UserStorer.Update(user)
}

func (us UserService) Delete(id int) (*user.User, error) {
	return us.UserStorer.Delete(id)
}
