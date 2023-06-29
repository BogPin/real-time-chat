package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/BogPin/real-time-chat/backend/api/models"
	"github.com/BogPin/real-time-chat/backend/api/utils"
)

type IChatService interface {
	Create(userId int, chat models.ChatFromRequest) (*models.Chat, utils.HttpError)
	GetOne(userId, chatId int) (*models.Chat, utils.HttpError)
	GetUserChats(userId int) ([]models.Chat, utils.HttpError)
	Update(userId int, chat models.Chat) (*models.Chat, utils.HttpError)
	Delete(userId, chatId int) (*models.Chat, utils.HttpError)
}

type ChatService struct {
	ChatStorer         models.IChatStorer
	ParticipantStorer  models.IParticipantStorer
	MessageStorer      models.IMessageStorer
	ParticipantService IParticipantService
}

func NewChatService(chatStorer models.IChatStorer, participantStorer models.IParticipantStorer, messageStorer models.IMessageStorer, participantService ParticipantService) ChatService {
	return ChatService{
		ChatStorer:         chatStorer,
		ParticipantStorer:  participantStorer,
		MessageStorer:      messageStorer,
		ParticipantService: participantService,
	}
}

func (cs ChatService) Create(userId int, chatWithTitle models.ChatFromRequest) (*models.Chat, utils.HttpError) {
	dto := models.ChatDTO{Title: chatWithTitle.Title, CreatorId: userId}

	tx, err := cs.ChatStorer.Begin()
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				log.Println(err)
			}
		} else {
			err := tx.Commit()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	chat, err := cs.ChatStorer.CreateInTx(tx, dto)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	admin := models.Participant{UserId: userId, ChatId: chat.Id, Role: "admin"}
	_, err = cs.ParticipantStorer.CreateInTx(tx, admin)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	return chat, nil
}

func (cs ChatService) GetOne(userId, chatId int) (*models.Chat, utils.HttpError) {
	chat, err := cs.ChatStorer.GetOne(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := fmt.Errorf("no chat with id %d", chatId)
			return nil, utils.NewHttpError(err, http.StatusNotFound)
		}
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	userInChat, err := cs.ParticipantService.UserInChat(userId, chatId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}

	if !userInChat {
		err := fmt.Errorf("user %d doesn't participate in chat %d", userId, chatId)
		return nil, utils.NewHttpError(err, http.StatusForbidden)
	}

	return chat, nil
}

func (cs ChatService) GetUserChats(userId int) ([]models.Chat, utils.HttpError) {
	chats, err := cs.ChatStorer.GetUserChats(userId)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError)
	}
	return chats, nil
}

func (cs ChatService) Update(userId int, chat models.Chat) (*models.Chat, utils.HttpError) {
	userInChat, err := cs.ParticipantService.UserInChat(userId, chat.Id)
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

func (cs ChatService) Delete(userId, chatId int) (*models.Chat, utils.HttpError) {
	userInChat, err := cs.ParticipantService.UserInChat(userId, chatId)
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
