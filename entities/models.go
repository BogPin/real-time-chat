package entities

import "database/sql"

type DbModel struct {
	conn *sql.DB
}
