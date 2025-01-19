package integration_test_test

import (
	"QYRGYN/config"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"QYRGYN/controllers"
	"QYRGYN/database"
	"QYRGYN/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostIntegration(t *testing.T) {

	gin.SetMode(gin.TestMode)

	database.InitTestDatabase(config.GetDatabaseURL())
	defer database.CloseTestDatabase()

	router := gin.Default()
	router.POST("/posts", controllers.CreatePost)

	user := models.User{Username: "test_user", Email: "test@example.com"}
	database.DB.Create(&user)

	payload := []byte(`content=This+is+a+test.&userId=14`)

	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", user.ID)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)

	var post models.Post
	err := database.DB.First(&post).Error
	assert.NoError(t, err)
	assert.Equal(t, "This is a test.", post.Content)
	assert.Equal(t, int(user.ID), post.UserId)
}
