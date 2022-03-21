package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func GetDatabase() (*Database, error) {
	// PSQL connection
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGDATABASE"), os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error while connecting PSQL Database: ", err)
	}

	// Config
	err = godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error while loading .env file: %s", err)
	}

	return &Database{
		Db:   db,
		Sqlx: sqlx.NewDb(db, "postgresql"),
	}, nil
}
