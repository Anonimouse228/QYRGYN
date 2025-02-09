package models

// PaymentResponse структура для ответа от микросервиса оплаты
type PaymentResponse struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id"`
	Message       string `json:"message"`
}
