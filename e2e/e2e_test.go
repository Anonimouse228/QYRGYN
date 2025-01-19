package e2e

import (
	"QYRGYN/config"
	"QYRGYN/database"
	"QYRGYN/middleware"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	database.InitDatabase(config.GetDatabaseURL()) // Initialize the test database

	router := gin.Default()

	// Add middleware
	middleware.SetupLogger()
	router.Use(middleware.Logger())

	// Define routes (replace InitRoutes with your actual route setup function)
	router.GET("/posts", func(c *gin.Context) {
		// Simulate returning posts
		c.JSON(http.StatusOK, gin.H{
			"posts": []string{"Test Post Title", "Another Test Post"},
		})
	})

	return router
}

// E2E Test: User Login and Create Post
func TestE2E_UserLoginAndCreatePost(t *testing.T) {
	// Initialize the application (main.go should set up the server and routes)
	router := setupTestRouter()

	// Step 1: Simulate User Registration
	regBody := map[string]string{
		"username": "testuser",
		"email":    "testuser@example.com",
		"password": "password123",
	}
	regBodyJSON, _ := json.Marshal(regBody)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(regBodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Step 2: Simulate User Login
	loginBody := map[string]string{
		"email":    "testuser@example.com",
		"password": "password123",
	}
	loginBodyJSON, _ := json.Marshal(loginBody)

	req, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginBodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Extract token from login response
	var loginResponse map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]
	assert.NotEmpty(t, token)

	// Step 3: Simulate Post Creation
	postBody := map[string]string{
		"title":   "My E2E Test Post",
		"content": "This is the content of the test post.",
	}
	postBodyJSON, _ := json.Marshal(postBody)

	req, _ = http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(postBodyJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // Pass token in Authorization header
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Step 4: Verify Post Retrieval
	req, _ = http.NewRequest(http.MethodGet, "/posts", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "My E2E Test Post")
}
