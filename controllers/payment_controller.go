package controllers

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// PaymentPage renders the payment page for the user to purchase a subscription.
func PaymentPage(c *gin.Context) {
	// Retrieve the user ID and the subscription to be purchased
	subscriptionID := c.Param("subscription_id")
	userID := c.GetUint("userID")

	var subscription models.Subscription
	if err := database.DB.First(&subscription, subscriptionID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription not found"})
		return
	}

	amount := subscription.Price
	c.HTML(http.StatusOK, "payment.html", gin.H{
		"userID":         userID,
		"subscriptionID": subscriptionID,
		"amount":         amount,
	})
}

// Payment handles the form submission for a user paying for a subscription.
func Payment(c *gin.Context) {
	// Получаем userID из контекста
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var form struct {
		CardName       string  `form:"card_name" binding:"required"`
		CardNumber     string  `form:"card_number" binding:"required,len=16"`
		ExpiryDate     string  `form:"expiry_date" binding:"required"`
		CVV            string  `form:"cvv" binding:"required,len=3"`
		SubscriptionID uint    `form:"subscription_id" binding:"required"`
		Amount         float64 `form:"amount" binding:"required"`
	}

	// Привязка данных формы
	if err := c.ShouldBind(&form); err != nil {
		// Печатаем подробности ошибки для отладки
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid payment details",
			"details": err.Error(),
		})
		return
	}

	// Имитация задержки при обработке платежа
	time.Sleep(2 * time.Second)

	// Извлекаем подписку из базы
	var subscription models.Subscription
	if err := database.DB.First(&subscription, form.SubscriptionID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription not found"})
		return
	}

	// Проверяем, была ли уже оплачена подписка
	if subscription.Status == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription already paid"})
		return
	}

	// Обновляем статус подписки на "paid"
	subscription.Status = "paid"
	if err := database.DB.Save(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription status"})
		return
	}

	// Возвращаем успех или перенаправляем на страницу подписки
	c.JSON(http.StatusOK, gin.H{"message": "Payment successful", "subscription": subscription})
}

// ProcessPayment processes the payment and updates the subscription status.
func ProcessPayment(c *gin.Context) {
	userID, _ := c.Get("userID")
	var request struct {
		CardName       string  `form:"card_name"`
		CardNumber     string  `form:"card_number"`
		ExpiryDate     string  `form:"expiry_date"`
		CVV            string  `form:"cvv"`
		SubscriptionID uint    `form:"subscription_id"`
		Amount         float64 `form:"amount"`
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var subscription models.Subscription
	if err := database.DB.First(&subscription, request.SubscriptionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	// Simulate payment processing (for the sake of the example)
	time.Sleep(2 * time.Second)

	// If payment is successful, update subscription status
	if subscription.Status != "paid" {
		subscription.Status = "paid"
		database.DB.Save(&subscription)
		c.JSON(http.StatusOK, gin.H{"message": "Subscription activated", "subscription": subscription})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription already paid"})
	}
}
