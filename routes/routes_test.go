package routes

import (
	"QYRGYN/config"
	"QYRGYN/database"
	"QYRGYN/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper function to set up routes for testing
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

// Test for the GET /posts route
func TestGetPosts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize router
	router := setupTestRouter()

	// Create HTTP request
	req, _ := http.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	// Send request to the router
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Post Title")
}
