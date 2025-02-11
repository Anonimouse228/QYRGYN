package models

type PaymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
