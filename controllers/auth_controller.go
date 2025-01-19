package controllers

import (
	"QYRGYN/config"
	"QYRGYN/database"
	"QYRGYN/models"
	"QYRGYN/util"
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var store = sessions.NewCookieStore([]byte(config.GetSecretKey()))

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func Register(c *gin.Context) {
	// Get form data
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	println(username)
	println(email)
	println(password)
	// Validate input
	if username == "" || email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "All fields are required."})
		return
	}

	var existingUser models.User
	database.DB.Where("email = ?", email).First(&existingUser)
	if existingUser.ID != 0 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Failed to register user."})
		return
	}

	// Create email verification token
	token, err := generateToken()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Token generation failed"})
		return
	}

	// Create user
	user := models.User{
		Username:          username,
		Email:             email,
		Password:          string(hashedPassword),
		VerificationToken: token,
		Verified:          false,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Failed to register user."})
		return
	}

	// Send verification email
	err = util.SendVerificationEmail(user.Email, token)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to send verification email"})
		return
	}

	// Redirect to log in
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
	session.Values["role"] = user.Role

	err := session.Save(c.Request, c.Writer)
	if err != nil {
		println("SESSIONERERERERERERERR", err.Error())
		return
	}

	if user.Role == "admin" {
		// Redirect to admin panel
		c.Redirect(http.StatusFound, "/admin/posts")
	} else {
		// Redirect to posts
		c.Redirect(http.StatusFound, "/posts")
	}
}

func Logout(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	delete(session.Values, "userID")
	err := session.Save(c.Request, c.Writer)
	if err != nil {
		log.Fatal(err)
		return ///////////////////ТУТ ХЗ ЧТО СДЕЛАЛ ААААААААААААААААААААААААААААААААААААА
	}

	c.Redirect(http.StatusFound, "/login")
}

func RegisterHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func LoginHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}
