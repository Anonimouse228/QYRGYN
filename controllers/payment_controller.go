package controllers

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

// PaymentPage renders the payment page for the user to purchase a subscription.
func PaymentPage(c *gin.Context) {
	// Retrieve the user ID and the subscription to be purchased
	//subscriptionID := c.Param("subscription_id")
	userID := c.GetUint("userID")
	//
	//var subscription models.Subscription
	//if err := database.DB.First(&subscription, subscriptionID).Error; err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription not found"})
	//	return
	//}

	//amount := subscription.Price
	c.HTML(http.StatusOK, "payment.html", gin.H{
		"userID": userID,
		//"subscriptionID": subscriptionID,
		"amount": "9.99$",
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
		CardName     string  `form:"cardName" binding:"required"`
		CardNumber   string  `form:"cardNumber" binding:"required,len=16"`
		ExpiryDate   string  `form:"expiryDate" binding:"required"`
		CVV          string  `form:"cvv" binding:"required,len=3"`
		Subscription uint    `form:"subscription" binding:"required"`
		Amount       float64 `form:"amount" binding:"required"`
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

	// Имитация задержки при обработке платежа (удалить или заменить)
	time.Sleep(2 * time.Second)

	// Извлекаем подписку из базы
	var subscription models.Subscription
	if err := database.DB.First(&subscription, form.Subscription).Error; err != nil {
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
//func ProcessPayment(c *gin.Context) {
//	userID, _ := c.Get("userID")
//	var paymentData struct {
//		CardName   string `form:"card_name"`
//		CardNumber string `form:"card_number"`
//		ExpiryDate string `form:"expiry_date"`
//		CVV        string `form:"cvv"`
//
//		Amount float64 `form:"amount"`
//	}
//
//	if err := c.ShouldBind(&paymentData); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
//		return
//	}
//	paymentData.Amount = 9.99
//
//	var user models.User
//	if err := database.DB.First(&user, userID).Error; err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//		return
//	}
//
//	jsonData, _ := json.Marshal(paymentData)
//	resp, err := http.Post("http://localhost:8080/api/payment/process", "application/json", bytes.NewBuffer(jsonData))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment service unavailable"})
//		return
//	}
//	defer resp.Body.Close()
//
//	// Simulate payment processing (for the sake of the example)
//	time.Sleep(2 * time.Second)
//
//	// If payment is successful, update subscription status
//	if subscription.Status != "paid" {
//		subscription.Status = "paid"
//		database.DB.Save(&subscription)
//		c.JSON(http.StatusOK, gin.H{"message": "Subscription activated", "subscription": subscription})
//	} else {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription already paid"})
//	}
//}

func ProcessPayment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		//c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		//return
	}

	// Преобразуем userID к uint
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Получаем данные карты
	var paymentData models.PaymentRequest
	if err := c.ShouldBindJSON(&paymentData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	fmt.Println("Raw request body:", string(body))

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Создаем подписку в базе данных
	subscription := models.Subscription{
		UserID: userIDUint,
		Type:   "premium",
		Price:  9.99,
		Status: "pending",
	}

	if err := database.DB.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	paymentData.Amount = subscription.Price
	paymentData.UserID = userIDUint
	paymentData.SubscriptionID = subscription.ID
	paymentData.Email = user.Email

	println(paymentData.CardName)
	println(paymentData.SubscriptionID)
	println(paymentData.ExpiryDate)
	println(paymentData.CVV)
	println(paymentData.CardNumber)

	// Отправляем запрос в платежный микросервис
	jsonData, err := json.Marshal(paymentData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare request"})
		return
	}

	resp, err := http.Post("https://qyrgyn-microservice.onrender.com/api/payments/process", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment service unavailable"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	fmt.Println("Raw response:", string(respBody)) // Посмотрим, что реально приходит

	var paymentResponse models.PaymentResponse
	if err := json.Unmarshal(respBody, &paymentResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from payment service"})
		return
	}

	println(paymentResponse.Success)
	println(paymentResponse.Message)

	if err := database.DB.First(&user, userIDUint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if paymentResponse.Success {
		subscription.Status = "paid"
		user.HasPremium = true
		database.DB.Save(&user)
	} else {
		subscription.Status = "failed"
	}

	database.DB.Save(&subscription)

	// Возвращаем ответ
	if paymentResponse.Success {
		c.JSON(http.StatusOK, gin.H{"message": "Subscription activated"})
	} else {
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "Payment failed", "message": paymentResponse.Message})
	}
}
