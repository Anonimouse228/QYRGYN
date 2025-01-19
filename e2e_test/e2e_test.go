package e2e_test

import (
	"QYRGYN/config"
	"QYRGYN/database"
	"QYRGYN/main_directory"
	"QYRGYN/models"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePostE2E(t *testing.T) {
	// Set up the application
	router := main_directory.SetupRouter()

	// Initialize a test database (in-memory or test DB)
	database.InitTestDatabase(config.GetDatabaseURL()) // Create this function for your test database initialization
	defer database.CloseTestDatabase()

	// Create a test user
	user := models.User{Username: "test_user", Email: "test@example.com"}
	database.DB.Create(&user)

	// Add middleware to simulate a logged-in user (set userID in context)
	router.Use(func(c *gin.Context) {
		c.Set("userID", user.ID)
		c.Next()
	})

	// Create a test payload
	payload := []byte(`content=This+is+a+test.`)

	// Create a request
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusFound, w.Code) // Check for redirect status

	// Verify the post in the database
	var post models.Post
	err := database.DB.First(&post).Error
	assert.NoError(t, err)

	assert.Equal(t, "This is a test.", post.Content)
	assert.Equal(t, int(user.ID), post.UserId)
}
