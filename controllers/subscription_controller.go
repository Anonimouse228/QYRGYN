package controllers

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Subscribe создает подписку для пользователя
func Subscribe(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var form struct {
		Type     string  `form:"type" binding:"required"`
		Price    float64 `form:"price" binding:"required"`
		Duration int     `form:"duration" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription details"})
		return
	}

	subscription := models.Subscription{
		UserID: userID,
		Type:   form.Type,
		Price:  form.Price,
		Status: "active",
	}

	if err := database.DB.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription created successfully"})
}

// GetSubscription получает подписку пользователя
func GetSubscription(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var subscription models.Subscription
	if err := database.DB.Where("user_id = ?", userID).First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active subscription found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// CancelSubscription отменяет подписку пользователя
func CancelSubscription(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := database.DB.Where("user_id = ?", userID).Delete(&models.Subscription{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription canceled successfully"})
}
