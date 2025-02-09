package database

import (
	"QYRGYN/models"
	"log"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDatabase(connectionString string) {
	var err error
	DB, err = gorm.Open("postgres", connectionString)

	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Chat{}, &models.Message{}).Error
	if err != nil {
		log.Fatalf("Could not migrate the database: %v", err)
	}
}
