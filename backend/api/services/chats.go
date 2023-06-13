package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/api/models/chat"
	"github.com/BogPin/real-time-chat/backend/api/models/message"
	"github.com/BogPin/real-time-chat/backend/api/models/participant"
	"github.com/BogPin/real-time-chat/backend/api/utils"
)

type Chat interface {
	Create(userId int, chat chat.ChatFromRequest) (*chat.Chat, utils.HttpError)
	GetOne(userId, chatId int) (*chat.Chat, utils.HttpError)
	GetUserChats(userId int) ([]chat.Chat, utils.HttpError)
	Update(userId int, chat chat.Chat) (*chat.Chat, utils.HttpError)
	Delete(userId, chatId int) (*chat.Chat, utils.HttpError)
}

type ChatService struct {
	ChatStorer        *chat.ChatStorer
	ParticipantStorer *participant.ParticipantStorer
	MessageStorer     *message.MessageStorer
}

func (cs ChatService) Create(userId int, chatWithTitle chat.ChatFromRequest) (*chat.Chat, utils.HttpError) {
	dto := chat.ChatDTO{Title: chatWithTitle.Title, CreatorId: userId}

	tx, err := cs.ChatStorer.DB.Begin()
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	chat, err := cs.ChatStorer.CreateInTx(tx, dto)
	if err != nil {
		tx.Rollback()
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	admin := participant.Participant{UserId: userId, ChatId: chat.Id, Role: "admin"}
	_, err = cs.ParticipantStorer.CreateInTx(tx, admin)
	if err != nil {
		tx.Rollback()
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	tx.Commit()
	return chat, nil
}

func (cs ChatService) GetOne(userId, chatId int) (*chat.Chat, utils.HttpError) {
	userInChat, err := userInChat(cs.ChatStorer.DB, userId, chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if !userInChat {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	chat, err := cs.ChatStorer.GetOne(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := fmt.Errorf("no chat with id %d", chatId)
			return nil, utils.NewHttpError(err, http.StatusNotFound)
		}
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}
	return chat, nil
}

func (cs ChatService) GetUserChats(userId int) ([]chat.Chat, utils.HttpError) {
	chats, err := cs.ChatStorer.GetUserChats(userId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}
	return chats, nil
}

func (cs ChatService) Update(userId int, chat chat.Chat) (*chat.Chat, utils.HttpError) {
	userInChat, err := userInChat(cs.ChatStorer.DB, userId, chat.Id)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if !userInChat {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, chat.Id)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	user, err := cs.ParticipantStorer.GetOne(userId, chat.Id)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if user.Role != "admin" {
		err := fmt.Errorf("user %d doesn't have permission to update chat %d", userId, chat.Id)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	updChat, err := cs.ChatStorer.Update(chat)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}
	return updChat, nil
}

func (cs ChatService) Delete(userId, chatId int) (*chat.Chat, utils.HttpError) {
	userInChat, err := userInChat(cs.ChatStorer.DB, userId, chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if !userInChat {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	user, err := cs.ParticipantStorer.GetOne(userId, chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if user.Role != "admin" {
		err := fmt.Errorf("user %d doesn't have permission to update chat %d", userId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	_, err = cs.MessageStorer.DeleteAll(chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	_, err = cs.ParticipantStorer.DeleteAll(chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	dltChat, err := cs.ChatStorer.Delete(chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}
	return dltChat, nil
}
