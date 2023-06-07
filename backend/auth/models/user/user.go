package user

import (
	"database/sql"
)

type User struct {
	Id           int    `json:"id"`
	Tag          string `json:"tag"`
	PasswordHash string `json:"passwordHash"`
}

type Credentials struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

type Storer struct {
	DB *sql.DB
}

func (s Storer) GetUser(username string) (User, error) {
	var usr User
	query := "select id, tag, password from users where tag=$1"
	row := s.DB.QueryRow(query, username)
	err := row.Scan(&usr.Id, &usr.Tag, &usr.PasswordHash)
	return usr, err
}

func (s Storer) CreateUser(creds Credentials) (User, error) {
	var usr User
	query := "insert into users (tag, password) value ($1, $2) returning id, tag, password"
	row := s.DB.QueryRow(query, creds.Username, creds.PasswordHash)
	err := row.Scan(&usr.Id, &usr.Tag, &usr.PasswordHash)
	return usr, err
}
