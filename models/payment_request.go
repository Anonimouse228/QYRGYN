package models

// PaymentRequest структура для запроса на оплату
type PaymentRequest struct {
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	CardNumber string  `json:"card_number"`
	CardExpiry string  `json:"card_expiry"`
	CardCVV    string  `json:"card_cvv"`
	Name       string  `json:"name"`    // Имя владельца карты
	Address    string  `json:"address"` // Адрес владельца карты
}
