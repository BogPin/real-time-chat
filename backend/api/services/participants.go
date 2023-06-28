package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/api/models"
	"github.com/BogPin/real-time-chat/backend/api/utils"
	"golang.org/x/exp/slices"
)

type IParticipantService interface {
	Create(userId int, participant models.Participant) (*models.Participant, utils.HttpError)
	GetChatUsers(userId, chatId int) ([]models.ChatUser, utils.HttpError)
	Update(userId int, participant models.Participant) (*models.Participant, utils.HttpError)
	Delete(userId int, participant models.Participant) (*models.Participant, utils.HttpError)
	UserInChat(userId, chatId int) (bool, error)
}

type ParticipantService struct {
	ParticipantStorer models.IParticipantStorer
}

func NewParticipantService(participantStorer models.IParticipantStorer) ParticipantService {
	return ParticipantService{ParticipantStorer: participantStorer}
}

func (ps ParticipantService) Create(userId int, participant models.Participant) (*models.Participant, utils.HttpError) {
	chatId := participant.ChatId
	userInChat, err := ps.UserInChat(userId, chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if !userInChat {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	_, err = ps.ParticipantStorer.Create(participant)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	return &participant, nil
}

func (ps ParticipantService) GetChatUsers(userId, chatId int) ([]models.ChatUser, utils.HttpError) {
	chatUsers, err := ps.ParticipantStorer.GetChatUsers(chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	idx := slices.IndexFunc(chatUsers, func(chatUser models.ChatUser) bool { return chatUser.UserId == userId })
	if idx == -1 {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	return chatUsers, nil
}

func (ps ParticipantService) Update(userId int, participant models.Participant) (*models.Participant, utils.HttpError) {
	chatId := participant.ChatId
	userInChat, err := ps.UserInChat(userId, chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if !userInChat {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	user, err := ps.ParticipantStorer.GetOne(userId, chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if user.Role != "admin" {
		err := fmt.Errorf("user %d doesn't have permission to update user %d in chat %d", userId, participant.UserId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	newParticipant, err := ps.ParticipantStorer.Update(participant)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	return newParticipant, nil
}

func (ps ParticipantService) Delete(userId int, participant models.Participant) (*models.Participant, utils.HttpError) {
	chatId := participant.ChatId
	userInChat, err := ps.UserInChat(userId, chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if !userInChat {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	user, err := ps.ParticipantStorer.GetOne(userId, chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if user.Role != "admin" {
		err := fmt.Errorf("user %d doesn't have permission to delete user %d in chat %d", userId, participant.UserId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	dltParticipant, err := ps.ParticipantStorer.Delete(participant)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	return dltParticipant, nil
}

func (ps ParticipantService) UserInChat(userId, chatId int) (bool, error) {
	_, err := ps.ParticipantStorer.GetOne(userId, chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
