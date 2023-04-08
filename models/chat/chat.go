package chat

import "database/sql"

type ChatStorer struct {
	DB *sql.DB
}

type Chat struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	CreatorId int    `json:"creator_id"`
	CreatedAt string `json:"created_at"`
}

type ChatDTO struct {
	Title     string `json:"title"`
	CreatorId int    `json:"creator_id"`
	CreatedAt string `json:"created_at"`
}

func (cs *ChatStorer) GetAllChats() ([]Chat, error) {
	chats := make([]Chat, 0)
	query := "SELECT * FROM chats"
	rows, err := cs.DB.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var chat Chat
		err := rows.Scan(&chat.Id, &chat.Title, &chat.CreatorId, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (cs *ChatStorer) GetChat(id int) (*Chat, error) {
	var chat Chat
	query := "SELECT * FROM chats WHERE id = $1"
	row := cs.DB.QueryRow(query, id)
	err := row.Scan(&chat.Id, &chat.Title, &chat.CreatorId, &chat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &chat, err
}
func (cs *ChatStorer) Update(chat Chat) (*Chat, error) {
	var updChat Chat
	query := "UPDATE chats SET title=$1 WHERE id=$2 RETURNING *"
	row := cs.DB.QueryRow(query, chat.Title, chat.Id)
	err := row.Scan(&updChat.Id, &updChat.Title, &updChat.CreatorId, &updChat.CreatedAt)
	return &chat, err
}

func (cs *ChatStorer) Create(tdo ChatDTO) (*Chat, error) {
	var chat Chat
	query := "INSERT INTO chats (title, creator_id, created_at) VALUES ($1, $2, $3)"
	row := cs.DB.QueryRow(query, tdo.Title, tdo.CreatorId, tdo.CreatedAt)
	err := row.Scan(&chat.Id, &chat.Title, &chat.CreatorId, &chat.CreatedAt)
	return &chat, err
}

func (cs *ChatStorer) Delete(id int) (*Chat, error) {
	var chat Chat
	query := "DELETE FROM chats WHERE id = $1"
	row := cs.DB.QueryRow(query, id)
	err := row.Scan(&chat.Id, &chat.Title, &chat.CreatorId, &chat.CreatedAt)
	return &chat, err
}
