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

	DB.AutoMigrate(&models.User{}, &models.Post{}, &models.PaymentRequest{}, &models.PaymentResponse{}, &models.Subscription{})
}
