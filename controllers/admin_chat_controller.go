package controllers

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AdminChatList(c *gin.Context) {
	var chats []models.Chat

	database.DB.Where("status = ?", "active").Find(&chats)

	c.HTML(http.StatusOK, "admin_chats.html", gin.H{
		"chats": chats,
	})
}

func AdminChat(c *gin.Context) {
	chatID := c.Param("chatID")
	var chat models.Chat

	// Ensure chat exists
	if err := database.DB.Where("id = ?", chatID).First(&chat).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Chat not found"})
		return
	}

	// Fetch chat messages
	var messages []models.Message
	database.DB.Where("chat_id = ?", chatID).Order("created_at ASC").Find(&messages)

	// Render chat page
	c.HTML(http.StatusOK, "admin_chat.html", gin.H{
		"chatID":   chat.ID,
		"messages": messages,
	})
}

func AdminCloseChat(c *gin.Context) {
	chatID := c.Param("chatID")
	var chat models.Chat

	if err := database.DB.Where("id = ?", chatID).First(&chat).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Chat not found"})
		return
	}

	chat.Status = "inactive"
	chat.ClosedAt = time.Now()

	if err := database.DB.Save(&chat).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to update chat status"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/chats/")
}

func AdminSendMessage(c *gin.Context) {
	chatID := c.Param("chatID")
	messageContent := c.PostForm("message")

	var chat models.Chat
	if err := database.DB.Where("id = ?", chatID).First(&chat).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Chat not found"})
		return
	}

	message := models.Message{
		ChatID:   chatID,
		SenderID: "Admin",
		Content:  messageContent,
	}
	database.DB.Create(&message)

	c.Redirect(http.StatusSeeOther, "/admin/chat/"+chatID)
}
