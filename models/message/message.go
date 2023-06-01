package message

import "database/sql"

type MessageStorer struct {
	DB *sql.DB
}

type Message struct {
	Id       int    `json:"id"`
	SenderId int    `json:"sender_id"`
	ChatId   int    `json:"chat_id"`
	Type     string `json:"type"`
	Text     string `json:"text"`
}

type MessageDTO struct {
	SenderId int    `json:"sender_id"`
	ChatId   int    `json:"chat_id"`
	Type     string `json:"type"`
	Text     string `json:"text"`
}

func (cs *MessageStorer) GetAllMessages() ([]Message, error) {
	messages := make([]Message, 0)
	query := "SELECT * FROM messages"
	rows, err := cs.DB.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Id, &message.SenderId, &message.ChatId, &message.Type, &message.Text)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (cs *MessageStorer) GetMessage(id int) (*Message, error) {
	var message Message
	query := "SELECT * FROM messages WHERE id = $1"
	row := cs.DB.QueryRow(query, id)
	err := row.Scan(&message.Id, &message.SenderId, &message.ChatId, &message.Type, &message.Text)
	if err != nil {
		return nil, err
	}
	return &message, err
}
func (cs *MessageStorer) Update(message Message) (*Message, error) {
	var updMessage Message
	query := "UPDATE messages SET title=$1 WHERE id=$2 RETURNING *"
	row := cs.DB.QueryRow(query, message.Text, message.Id)
	err := row.Scan(&updMessage.Id, &updMessage.SenderId, &updMessage.ChatId, &updMessage.Type, &updMessage.Text)
	return &message, err
}

func (cs *MessageStorer) Create(tdo MessageDTO) (*Message, error) {
	var message Message
	query := "INSERT INTO messages (title, creator_id, created_at) VALUES ($1, $2, $3)"
	row := cs.DB.QueryRow(query, tdo.Text, tdo.SenderId, tdo.ChatId, tdo.Type)
	err := row.Scan(&message.Id, &message.SenderId, &message.ChatId, &message.Type, &message.Text)
	return &message, err
}

func (cs *MessageStorer) Delete(id int) (*Message, error) {
	var message Message
	query := "DELETE FROM messages WHERE id = $1"
	row := cs.DB.QueryRow(query, id)
	err := row.Scan(&message.Id, &message.SenderId, &message.ChatId, &message.Type, &message.Text)
	return &message, err
}
