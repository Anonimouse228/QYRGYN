package controllers

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func Register(c *gin.Context) {
	// Get form data
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validate input
	if username == "" || email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "All fields are required."})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Failed to register user."})
		return
	}

	// Create user
	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Failed to register user."})
		return
	}

	// Redirect to login
	c.Redirect(http.StatusFound, "/login")
}

func Login(c *gin.Context) {
	// Get form data
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validate input
	if email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Email and password are required."})
		return
	}

	// Find user
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid email or password."})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid email or password."})
		return
	}

	// Create session
	session, _ := store.Get(c.Request, "session")
	session.Values["userID"] = user.ID
	session.Save(c.Request, c.Writer)

	// Redirect to posts
	c.Redirect(http.StatusFound, "/posts")
}

func Logout(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	delete(session.Values, "userID")
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/posts")
}

func RegisterHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func LoginHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}
