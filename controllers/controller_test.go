package controllers

import (
	"QYRGYN/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Fake database
type FakeDB struct {
	posts []models.Post
}

func (f *FakeDB) Create(post *models.Post) error {
	f.posts = append(f.posts, *post)
	return nil
}

func TestCreatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/posts", CreatePost)

	tests := []struct {
		name         string
		userID       interface{}
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Unauthorized user",
			userID:       nil,
			body:         `{"content":"Test Content"}`,
			expectedCode: http.StatusUnauthorized,
			expectedBody: "In case somehow the userID is not set in the context (shouldn't happen)",
		},
		{
			name:         "Binding error",
			userID:       uint(1),
			body:         `{"invalid":"data"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Binding error:",
		},
		{
			name:         "Successful post creation",
			userID:       uint(1),
			body:         `{"content":"Test Content"}`,
			expectedCode: http.StatusFound,
			expectedBody: "/posts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Set("userID", tt.userID)

			req, _ := http.NewRequest(http.MethodPost, "/posts", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}
