package models

import (
	"database/sql"
)

type Participant struct {
	UserId int    `json:"userId"`
	ChatId int    `json:"chatId"`
	Role   string `json:"role"`
}

type ChatUser struct {
	Participant
	Name string `json:"name"`
}

type ParticipantFromRequest struct {
	UserId int `json:"userId"`
	ChatId int `json:"chatId"`
}

type IParticipantStorer interface {
	Create(participant Participant) (*Participant, error)
	CreateInTx(tx *sql.Tx, participant Participant) (*Participant, error)
	GetOne(userId, chatId int) (*Participant, error)
	GetChatUsers(chatId int) ([]ChatUser, error)
	Update(participant Participant) (*Participant, error)
	Delete(participant Participant) (*Participant, error)
	DeleteAll(chatId int) (sql.Result, error)
}

type ParticipantStorer struct {
	DB *sql.DB
}

func NewParticipantStorer(db *sql.DB) ParticipantStorer {
	return ParticipantStorer{DB: db}
}

func (ps ParticipantStorer) Create(participant Participant) (*Participant, error) {
	query := "INSERT INTO participants (user_id, chat_id, role) VALUES ($1, $2, $3)"
	_, err := ps.DB.Exec(query, participant.UserId, participant.ChatId, participant.Role)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (ps ParticipantStorer) CreateInTx(tx *sql.Tx, participant Participant) (*Participant, error) {
	query := "INSERT INTO participants (user_id, chat_id, role) VALUES ($1, $2, $3)"
	_, err := tx.Exec(query, participant.UserId, participant.ChatId, participant.Role)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (ps ParticipantStorer) GetOne(userId, chatId int) (*Participant, error) {
	var participant Participant
	query := "SELECT * FROM participants WHERE user_id=$1 AND chat_id=$2"
	row := ps.DB.QueryRow(query, userId, chatId)
	err := row.Scan(&participant.UserId, &participant.ChatId, &participant.Role)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (ps ParticipantStorer) GetChatUsers(chatId int) ([]ChatUser, error) {
	chatUsers := make([]ChatUser, 0)
	query := "SELECT u.id, u.name, p.chat_id, p.role FROM participants p JOIN users u ON p.user_id=u.id WHERE p.chat_id=$2"
	rows, err := ps.DB.Query(query, chatId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var chatUser ChatUser
		err := rows.Scan(&chatUser.UserId, &chatUser.Name, &chatUser.ChatId, &chatUser.Role)
		if err != nil {
			return nil, err
		}
		chatUsers = append(chatUsers, chatUser)
	}
	return chatUsers, nil
}

func (ps ParticipantStorer) Update(participant Participant) (*Participant, error) {
	var updParticipant Participant
	query := "UPDATE participants SET role=$1 WHERE user_id=$1 AND chat_id=$2 RETURNING *"
	row := ps.DB.QueryRow(query, participant.Role, participant.UserId, participant.ChatId)
	err := row.Scan(&updParticipant.UserId, &updParticipant.ChatId, &updParticipant.Role)
	if err != nil {
		return nil, err
	}
	return &updParticipant, nil
}

func (ps ParticipantStorer) Delete(participant Participant) (*Participant, error) {
	var dltParticipant Participant
	query := "DELETE FROM participants WHERE user_id=$1 AND chat_id=$2 RETURNING *"
	row := ps.DB.QueryRow(query, participant.UserId, participant.ChatId)
	err := row.Scan(&dltParticipant.UserId, &dltParticipant.ChatId, &dltParticipant.Role)
	if err != nil {
		return nil, err
	}
	return &dltParticipant, nil
}

func (ps ParticipantStorer) DeleteAll(chatId int) (sql.Result, error) {
	query := "DELETE FROM participants WHERE chat_id=$1"
	return ps.DB.Exec(query, chatId)
}
