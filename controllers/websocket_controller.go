package controllers

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Message)

func HandleConnections(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	log.Println("New WebSocket connection established")

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			delete(clients, conn)
			break
		}

		msg.CreatedAt = time.Now()

		result := database.DB.Create(&msg)
		if result.Error != nil {
			log.Println("Ошибка сохранения в БД:", result.Error)
			continue
		}

		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				println("Error sending message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func StartChat(c *gin.Context) {
	id := c.Param("id")
	//sessionUserID := c.GetUint("userID")
	//println(sessionUserID)

	userID, err := strconv.Atoi(id)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid user ID"})
		return
	}

	var existingChat models.Chat
	err = database.DB.Where("user_id = ? AND status = 'active'", userID).First(&existingChat).Error

	if err == nil {
		println(existingChat.ID, existingChat.Status, existingChat.UserID)
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/chat/%d", existingChat.ID))
		return
	}

	newChat := models.Chat{
		UserID:    userID,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&newChat).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to start chat"})
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/chat/%d", newChat.ID))
}

func ChatPage(c *gin.Context) {
	chatID := c.Param("chatID")
	var chat models.Chat

	if err := database.DB.Where("id = ?", chatID).First(&chat).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Chat not found"})
		return
	}

	if chat.Status == "Inactive" {
		c.HTML(http.StatusForbidden, "error.html", gin.H{"error": "Chat is inactive"})
		return
	}

	var messages []models.Message
	database.DB.Where("chat_id = ?", chatID).Order("created_at ASC").Find(&messages)

	c.HTML(http.StatusOK, "chat.html", gin.H{
		"chatID":   chat.ID,
		"messages": messages,
	})
}

func SendMessage(c *gin.Context) {
	chatID := c.Param("chatID")
	messageContent := c.PostForm("message")

	var chat models.Chat
	if err := database.DB.Where("id = ?", chatID).First(&chat).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Chat not found"})
		return
	}
	print(chatID, chat.Status, messageContent)

	message := models.Message{
		ChatID:   chatID,
		SenderID: "User",
		Content:  messageContent,
	}
	database.DB.Create(&message)

	c.Redirect(http.StatusSeeOther, "/chat/"+chatID)
}
