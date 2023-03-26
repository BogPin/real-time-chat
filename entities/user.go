package entities

import (
	"fmt"
	_ "github.com/lib/pq"
)

type User struct {
	id          int
	tag         string
	name        string
	password    string
	description string
}

type UserDTO struct {
	tag         string `json:"tag"`
	name        string `json:"name"`
	password    string `json:"password"`
	description string `json:"description"`
}

func (m *DbModel) CreateUser(dto UserDTO) (User, error) {
	var user User
	query := fmt.Sprintf("INSERT INTO users (tag, name, password, description) VALUES ('%s', '%s', '%s', '%s') RETURNING *",
		dto.tag, dto.name, dto.password, dto.description)
	err := m.conn.QueryRow(query).Scan(&user.id, &user.tag, &user.name, &user.password, &user.description)
	return user, err
}
