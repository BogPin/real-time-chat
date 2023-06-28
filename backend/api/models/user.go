package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

//TODO: use "null" package to make some fields sql nullable(e.g. description)

type IUserStorer interface {
	GetOne(id int) (*User, error)
	Update(user User) (*User, error)
	Delete(id int) (*User, error)
}

type UserStorer struct {
	DB *sql.DB
}

func NewUserStorer(db *sql.DB) UserStorer {
	return UserStorer{DB: db}
}

type User struct {
	Id          int    `json:"id"`
	Tag         string `json:"tag"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

type UserDTO struct {
	Tag         string `json:"tag"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

func (us UserStorer) GetOne(id int) (*User, error) {
	var user User
	query := "SELECT * FROM users where id = $1"
	row := us.DB.QueryRow(query, id)
	err := row.Scan(&user.Id, &user.Tag, &user.Name, &user.Password, &user.Description)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us UserStorer) Update(user User) (*User, error) {
	var updUser User
	query := "UPDATE users SET tag=$1, name=$2, password=$3, description=$4 WHERE id=$5 RETURNING *"
	row := us.DB.QueryRow(query, user.Tag, user.Name, user.Password, user.Description, user.Id)
	err := row.Scan(&updUser.Id, &updUser.Tag, &updUser.Name, &updUser.Password, &updUser.Description)
	if err != nil {
		return nil, err
	}
	return &updUser, nil
}

func (us UserStorer) Delete(id int) (*User, error) {
	var dltUser User
	query := "DELETE FROM users WHERE id=$1 RETURNING *"
	row := us.DB.QueryRow(query, id)
	err := row.Scan(&dltUser.Id, &dltUser.Tag, &dltUser.Name, &dltUser.Password, &dltUser.Description)
	if err != nil {
		return nil, err
	}
	return &dltUser, nil
}
