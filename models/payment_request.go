package models

// PaymentRequest структура для запроса на оплату
type PaymentRequest struct {
	CardName       string  `json:"card_name"`
	CardNumber     string  `json:"card_number"`
	ExpiryDate     string  `json:"expiry_date"`
	CVV            string  `json:"cvv"`
	Amount         float64 `json:"amount"`
	UserID         uint    `json:"user_id"`
	SubscriptionID uint    `json:"subscription_id"`
	Email          string  `json:"email"`
}
