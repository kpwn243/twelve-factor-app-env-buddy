package internal

import (
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
)

var database *gorm.DB

func InitDbConnection() (*gorm.DB, error) {
	config := GetConfiguration()
	db, err := gorm.Open("sqlite3", config.DbFileLocation)
	if err != nil {
		return nil, err
	}
	database = db

	return db, nil
}

func GetDbConnection() *gorm.DB {
	return database
}
