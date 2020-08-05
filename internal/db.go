package internal

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func InitDbConnection() (*sql.DB, error) {
	config := GetConfiguration()
	db, err := sql.Open("sqlite3", config.DbFileLocation)
	if err != nil {
		return nil, err
	}
	database = db

	return db, nil
}

func GetDbConnection() *sql.DB {
	return database
}
