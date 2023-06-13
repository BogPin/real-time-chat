package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/api/models/message"
	"github.com/BogPin/real-time-chat/backend/api/utils"
)

type Message interface {
	Create(userId int, MessageTDO message.MessageFromRequest) (*message.Message, utils.HttpError)
	GetOne(userId, messageId int) (*message.Message, utils.HttpError)
	GetChatMessages(userId, chatId, page int) ([]message.Message, utils.HttpError)
	Update(userId int, message message.Message) (*message.Message, utils.HttpError)
	Delete(userId int, message message.Message) (*message.Message, utils.HttpError)
}

type MessageService struct {
	MessageStorer *message.MessageStorer
}

func (ms MessageService) Create(userId int, fromRequest message.MessageFromRequest) (*message.Message, utils.HttpError) {
	chatId := fromRequest.ChatId
	userInChat, err := userInChat(ms.MessageStorer.DB, userId, chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if !userInChat {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	dto := message.MessageDTO{
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

func (ms MessageService) GetOne(userId, messageId int) (*message.Message, utils.HttpError) {
	msg, err := ms.MessageStorer.GetOne(messageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := fmt.Errorf("no message with id %d", messageId)
			return nil, utils.NewHttpError(err, http.StatusNotFound)
		}
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	userInChat, err := userInChat(ms.MessageStorer.DB, userId, msg.ChatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if !userInChat {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, msg.ChatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	return msg, nil
}

func (ms MessageService) GetChatMessages(userId, chatId, page int) ([]message.Message, utils.HttpError) {
	userInChat, err := userInChat(ms.MessageStorer.DB, userId, chatId)
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

func (ms MessageService) Update(userId int, message message.Message) (*message.Message, utils.HttpError) {
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

func (ms MessageService) Delete(userId int, message message.Message) (*message.Message, utils.HttpError) {
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
