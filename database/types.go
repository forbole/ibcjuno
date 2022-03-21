package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	Db   *sql.DB
	Sqlx *sqlx.DB
}
