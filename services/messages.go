package services

import (
	"github.com/BogPin/real-time-chat/models/message"
)

type Message interface {
	GetAllMessages() ([]message.Message, error)
	GetMessage(id int) (*message.Message, error)
	Update(message message.Message) (*message.Message, error)
	Delete(id int) (*message.Message, error)
	Create(MessageTDO message.MessageDTO) (*message.Message, error)
}

type MessageService struct {
	MessageStorer *message.MessageStorer
}

func (cs MessageService) GetAllMessages() ([]message.Message, error) {
	return cs.MessageStorer.GetAllMessages()
}

func (cs MessageService) GetMessage(id int) (*message.Message, error) {
	return cs.MessageStorer.GetMessage(id)
}

func (cs MessageService) Update(message message.Message) (*message.Message, error) {
	return cs.MessageStorer.Update(message)
}

func (cs MessageService) Create(messageDTO message.MessageDTO) (*message.Message, error) {
	return cs.MessageStorer.Create(messageDTO)
}

func (cs MessageService) Delete(id int) (*message.Message, error) {
	return cs.MessageStorer.Delete(id)
}
