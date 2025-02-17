package controllers

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"QYRGYN/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

// Serve the helpdesk page
func HelpdeskPageHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "helpdesk.html", nil)
}

// Handle helpdesk form submission
func HelpdeskController(c *gin.Context) {
	// Form data
	email := c.PostForm("email")
	subject := c.PostForm("subject")
	message := c.PostForm("message")

	// Validate input
	if email == "" || subject == "" || message == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "All fields are required"})
		return
	}

	// Handle attachments
	var attachments []string
	form, err := c.MultipartForm()
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Failed to parse form data"})
		return
	}
	files := form.File["attachments"]

	// Process uploaded files
	for _, file := range files {
		// Check file extension (basic security)
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".png" && ext != ".pdf" {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid file type"})
			return
		}

		// Save file
		filename := filepath.Join("uploads", file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "File upload failed"})
			return
		}
		attachments = append(attachments, filename)
	}

	// Compose email with sender's email included
	fullMessage := "From: " + email + "\n\n" + message

	// Send email
	err = util.SendEmail("suhansun13@gmail.com", subject, fullMessage, attachments)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to send email"})
		return
	}

	// Success response
	c.HTML(http.StatusOK, "error.html", gin.H{"error": "Email sent successfully!"})
}

func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Token is required"})
		return
	}

	var user models.User
	database.DB.Where("verification_token = ?", token).First(&user)
	if user.ID == 0 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid token"})
		return
	}

	if user.Verified {
		c.HTML(http.StatusOK, "error.html", gin.H{"error": "Email already verified"})
		return
	}

	// Update verification status
	user.Verified = true
	//user.VerificationToken = "" // Clear the token
	database.DB.Save(&user)

	c.HTML(http.StatusOK, "error.html", gin.H{"error": "Email verified successfully"})
}
