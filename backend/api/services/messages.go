package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/api/models"
	"github.com/BogPin/real-time-chat/backend/api/utils"
)

type IMessageService interface {
	Create(userId int, MessageTDO models.MessageFromRequest) (*models.Message, utils.HttpError)
	GetOne(userId, messageId int) (*models.Message, utils.HttpError)
	GetChatMessages(userId, chatId, page int) ([]models.Message, utils.HttpError)
	Update(userId int, message models.Message) (*models.Message, utils.HttpError)
	Delete(userId int, message models.Message) (*models.Message, utils.HttpError)
}

type MessageService struct {
	MessageStorer      models.IMessageStorer
	ParticipantService IParticipantService
}

func NewMessageService(messageStorer models.IMessageStorer, participantService IParticipantService) MessageService {
	return MessageService{
		MessageStorer:      messageStorer,
		ParticipantService: participantService,
	}
}

func (ms MessageService) Create(userId int, fromRequest models.MessageFromRequest) (*models.Message, utils.HttpError) {
	chatId := fromRequest.ChatId
	userInChat, err := ms.ParticipantService.UserInChat(userId, chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if !userInChat {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	dto := models.MessageDTO{
		SenderId: userId,
		ChatId:   fromRequest.ChatId,
		Type:     fromRequest.Type,
		Content:  fromRequest.Content,
	}

	msg, err := ms.MessageStorer.Create(dto)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	return msg, nil
}

func (ms MessageService) GetOne(userId, messageId int) (*models.Message, utils.HttpError) {
	msg, err := ms.MessageStorer.GetOne(messageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := fmt.Errorf("no message with id %d", messageId)
			return nil, utils.NewHttpError(err, http.StatusNotFound)
		}
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	userInChat, err := ms.ParticipantService.UserInChat(userId, msg.ChatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if !userInChat {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, msg.ChatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	return msg, nil
}

func (ms MessageService) GetChatMessages(userId, chatId, page int) ([]models.Message, utils.HttpError) {
	userInChat, err := ms.ParticipantService.UserInChat(userId, chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if !userInChat {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	messages, err := ms.MessageStorer.GetChatMessages(chatId, page)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	return messages, nil
}

func (ms MessageService) Update(userId int, message models.Message) (*models.Message, utils.HttpError) {
	originalMsg, httpErr := ms.GetOne(userId, message.Id)
	if httpErr != nil {
		return nil, httpErr
	}

	if userId != originalMsg.SenderId {
		err := fmt.Errorf("user %d didn't send message %d", userId, originalMsg.Id)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	msg, err := ms.MessageStorer.Update(message)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	return msg, nil
}

func (ms MessageService) Delete(userId int, message models.Message) (*models.Message, utils.HttpError) {
	originalMsg, httpErr := ms.GetOne(userId, message.Id)
	if httpErr != nil {
		return nil, httpErr
	}

	if userId != originalMsg.SenderId {
		err := fmt.Errorf("user %d didn't send message %d", userId, originalMsg.Id)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	msg, err := ms.MessageStorer.Delete(message.Id)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	return msg, nil
}
