package models

import (
	"database/sql"
)

type Chat struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	CreatorId int    `json:"creatorId"`
	CreatedAt string `json:"createdAt"`
}

type ChatDTO struct {
	Title     string `json:"title"`
	CreatorId int    `json:"creatorId"`
}

type ChatFromRequest struct {
	Title string `json:"title"`
}

type IChatStorer interface {
	Begin() (*sql.Tx, error)
	Create(dto ChatDTO) (*Chat, error)
	CreateInTx(tx *sql.Tx, dto ChatDTO) (*Chat, error)
	GetOne(id int) (*Chat, error)
	GetUserChats(userId int) ([]Chat, error)
	Update(chat Chat) (*Chat, error)
	Delete(id int) (*Chat, error)
}

type ChatStorer struct {
	DB *sql.DB
}

func NewChatStorer(db *sql.DB) ChatStorer {
	return ChatStorer{DB: db}
}

func (cs ChatStorer) Begin() (*sql.Tx, error) {
	return cs.DB.Begin()
}

func (cs ChatStorer) Create(dto ChatDTO) (*Chat, error) {
	var chat Chat
	query := "INSERT INTO chats (title, creator_id) VALUES ($1, $2) RETURNING *"
	row := cs.DB.QueryRow(query, dto.Title, dto.CreatorId)
	err := row.Scan(&chat.Id, &chat.Title, &chat.CreatorId, &chat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (cs ChatStorer) CreateInTx(tx *sql.Tx, dto ChatDTO) (*Chat, error) {
	var chat Chat
	query := "INSERT INTO chats (title, creator_id) VALUES ($1, $2) RETURNING *"
	row := tx.QueryRow(query, dto.Title, dto.CreatorId)
	err := row.Scan(&chat.Id, &chat.Title, &chat.CreatorId, &chat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (cs ChatStorer) GetOne(id int) (*Chat, error) {
	var chat Chat
	query := "SELECT * FROM chats WHERE id = $1"
	row := cs.DB.QueryRow(query, id)
	err := row.Scan(&chat.Id, &chat.Title, &chat.CreatorId, &chat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (cs ChatStorer) GetUserChats(userId int) ([]Chat, error) {
	userChats := make([]Chat, 0)
	query := "SELECT c.* FROM chats c JOIN participants p ON c.id = p.chat_id WHERE p.user_id = $1"
	rows, err := cs.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var chat Chat
		err := rows.Scan(&chat.Id, &chat.Title, &chat.CreatorId, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}
		userChats = append(userChats, chat)
	}
	return userChats, nil
}

func (cs ChatStorer) Update(chat Chat) (*Chat, error) {
	var updChat Chat
	query := "UPDATE chats SET title=$1 WHERE id=$2 RETURNING *"
	row := cs.DB.QueryRow(query, chat.Title, chat.Id)
	err := row.Scan(&updChat.Id, &updChat.Title, &updChat.CreatorId, &updChat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &updChat, nil
}

func (cs ChatStorer) Delete(id int) (*Chat, error) {
	var chat Chat
	query := "DELETE FROM chats WHERE id = $1 RETURNING *"
	row := cs.DB.QueryRow(query, id)
	err := row.Scan(&chat.Id, &chat.Title, &chat.CreatorId, &chat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}
