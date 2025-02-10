package models

// PaymentResponse структура для ответа от микросервиса оплаты
type PaymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
