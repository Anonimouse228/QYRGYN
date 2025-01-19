package database

import (
	"QYRGYN/models"
	_ "QYRGYN/models"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var test_DB *gorm.DB

func InitTestDatabase(connectionString string) {
	var err error
	test_DB, err = gorm.Open("postgres", connectionString)

	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	test_DB.AutoMigrate(&models.User{}, &models.Post{})
}
