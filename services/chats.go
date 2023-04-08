package services

import (
	"github.com/BogPin/real-time-chat/models/chat"
)

type Chat interface {
	GetAllChats() ([]chat.Chat, error)
	GetChat(id int) (*chat.Chat, error)
	Update(chat chat.Chat) (*chat.Chat, error)
	Delete(id int) (*chat.Chat, error)
	Create(ChatTDO chat.ChatDTO) (*chat.Chat, error)
}

type ChatService struct {
	ChatStorer *chat.ChatStorer
}

func (cs ChatService) GetAllChats() ([]chat.Chat, error) {
	return cs.ChatStorer.GetAllChats()
}

func (cs ChatService) GetChat(id int) (*chat.Chat, error) {
	return cs.ChatStorer.GetChat(id)
}

func (cs ChatService) Update(chat chat.Chat) (*chat.Chat, error) {
	return cs.ChatStorer.Update(chat)
}

func (cs ChatService) Create(chatDTO chat.ChatDTO) (*chat.Chat, error) {
	return cs.ChatStorer.Create(chatDTO)
}

func (cs ChatService) Delete(id int) (*chat.Chat, error) {
	return cs.ChatStorer.Delete(id)
}
