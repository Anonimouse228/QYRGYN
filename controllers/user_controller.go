package controllers

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	_ "reflect"
	"strconv"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "users.html", gin.H{"users": users})
}

func NewUserForm(c *gin.Context) {
	c.HTML(http.StatusOK, "new_user.html", nil)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := database.DB.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
		return
	}
	c.HTML(http.StatusOK, "user.html", gin.H{"user": user})
}

func EditUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := database.DB.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Post not found"})
		return
	}
	c.HTML(http.StatusOK, "edit_user.html", gin.H{"user": user})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := database.DB.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBind(&user); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&user)
	c.Redirect(http.StatusFound, "/users")
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("ID IS: ", id)
	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Could not delete user"})
		return
	}
	c.Redirect(http.StatusFound, "/users")
}

func GetUserProfile(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	// Get logged-in user ID from session
	sessionUserID := c.GetUint("userID")
	println(sessionUserID)
	// Fetch user details
	if err := database.DB.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
		return
	}

	// Render profile page with user info and session ID
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"user":          user,
		"sessionUserID": sessionUserID, // Pass session user ID to template
	})
}

func UpdateUserHTML(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := database.DB.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
		return
	}
	user.Password = ""
	c.HTML(http.StatusOK, "edit_user.html", gin.H{"user": user})
}

func UpdateUserProfile(c *gin.Context) {
	// Get user ID from URL parameter
	userID := c.Param("id")

	// Session validation
	sessionUserID := c.GetUint("userID") // Directly get session user ID as uint
	if sessionUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Convert userID to uint for comparison
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if uint(intUserID) != sessionUserID { // Check if user is editing their own profile
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}

	// Input validation
	var input struct {
		Username string `form:"username"`
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	// Bind form input
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	print(input.Username, "|", input.Email, "|", input.Password)
	// Check required fields
	if input.Username == "" || input.Email == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required"})
		return
	}

	// Validate email format
	//if !util.IsValidEmail(input.Email) {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
	//	return
	//}

	// Hash the password after validation
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update user profile
	result := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword), // Convert byte slice to string
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found or no changes made"})
		return
	}

	// Redirect to user profile page
	c.Redirect(http.StatusFound, "/users/"+userID)
}
