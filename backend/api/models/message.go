package models

import "database/sql"

type Message struct {
	Id        int    `json:"id"`
	SenderId  int    `json:"senderId"`
	ChatId    int    `json:"chatId"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type MessageDTO struct {
	SenderId int    `json:"senderId"`
	ChatId   int    `json:"chatId"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}

type MessageFromRequest struct {
	ChatId  int    `json:"chatId"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type IMessageStorer interface {
	Create(tdo MessageDTO) (*Message, error)
	GetOne(id int) (*Message, error)
	GetChatMessages(chatId, page int) ([]Message, error)
	Update(message Message) (*Message, error)
	Delete(id int) (*Message, error)
	DeleteAll(chatId int) (sql.Result, error)
}

const PAGE_SIZE = 50

type MessageStorer struct {
	DB *sql.DB
}

func NewMessageStorer(db *sql.DB) MessageStorer {
	return MessageStorer{DB: db}
}

func (cs MessageStorer) Create(tdo MessageDTO) (*Message, error) {
	var message Message
	query := "INSERT INTO messages (sender_id, chat_id, type, content) VALUES ($1, $2, $3, $4) RETURNING *"
	row := cs.DB.QueryRow(query, tdo.SenderId, tdo.ChatId, tdo.Type, tdo.Content)
	err := row.Scan(&message.Id, &message.SenderId, &message.ChatId, &message.Type, &message.Content, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (cs MessageStorer) GetOne(id int) (*Message, error) {
	var message Message
	query := "SELECT * FROM messages WHERE id = $1"
	row := cs.DB.QueryRow(query, id)
	err := row.Scan(&message.Id, &message.SenderId, &message.ChatId, &message.Type, &message.Content)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (cs MessageStorer) GetChatMessages(chatId, page int) ([]Message, error) {
	messages := make([]Message, 0)
	query := "SELECT * FROM messages WHERE chat_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"
	rows, err := cs.DB.Query(query, chatId, PAGE_SIZE, page*PAGE_SIZE)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Id, &message.SenderId, &message.ChatId, &message.Type, &message.Content)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (cs MessageStorer) Update(message Message) (*Message, error) {
	var updMessage Message
	query := "UPDATE messages SET content=$1 WHERE id=$2 RETURNING *"
	row := cs.DB.QueryRow(query, message.Content, message.Id)
	err := row.Scan(&updMessage.Id, &updMessage.SenderId, &updMessage.ChatId, &updMessage.Type, &updMessage.Content)
	if err != nil {
		return nil, err
	}
	return &updMessage, nil
}

func (cs MessageStorer) Delete(id int) (*Message, error) {
	var message Message
	query := "DELETE FROM messages WHERE id = $1"
	row := cs.DB.QueryRow(query, id)
	err := row.Scan(&message.Id, &message.SenderId, &message.ChatId, &message.Type, &message.Content)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (cs MessageStorer) DeleteAll(chatId int) (sql.Result, error) {
	query := "DELETE FROM messages WHERE chat_id = $1"
	return cs.DB.Exec(query, chatId)
}
