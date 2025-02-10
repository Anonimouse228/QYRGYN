package database

import (
	"QYRGYN/models"
	_ "QYRGYN/models"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func InitDatabase(connectionString string) {
	var err error
	DB, err = gorm.Open("postgres", connectionString)

	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Chat{}, &models.Message{}, &models.PaymentRequest{}, &models.PaymentResponse{}, &models.Subscription{}).Error
	if err != nil {
		log.Fatalf("Could not migrate the database: %v", err)
	}
}
