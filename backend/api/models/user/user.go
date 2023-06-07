package user

import (
	"database/sql"

	_ "github.com/lib/pq"
)

//TODO: use "null" package to make some fields sql nullable(e.g. description)

type UserStorer struct {
	DB *sql.DB
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

func (us *UserStorer) GetAll() ([]User, error) {
	users := make([]User, 0)
	query := "SELECT * FROM users"
	rows, err := us.DB.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Tag, &user.Name, &user.Password, &user.Description)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (us *UserStorer) GetOne(id int) (*User, error) {
	var user User
	query := "SELECT * FROM users where id = $1"
	row := us.DB.QueryRow(query, id)
	err := row.Scan(&user.Id, &user.Tag, &user.Name, &user.Password, &user.Description)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (us *UserStorer) Update(user User) (*User, error) {
	var newUser User
	query := "UPDATE users SET tag=$1, name=$2, password=$3, description=$4 WHERE id=$5 RETURNING *"
	row := us.DB.QueryRow(query, user.Tag, user.Name, user.Password, user.Description, user.Id)
	err := row.Scan(&newUser.Id, &newUser.Tag, &newUser.Name, &newUser.Password, &newUser.Description)
	return &user, err
}

func (us *UserStorer) Delete(id int) (*User, error) {
	var user User
	query := "DELETE FROM users WHERE id=$1 RETURNING *"
	row := us.DB.QueryRow(query, id)
	err := row.Scan(&user.Id, &user.Tag, &user.Name, &user.Password, &user.Description)
	return &user, err
}
