package models

import "time"

type Message struct {
	SenderID  string    `json:"sender_id"`
	ChatID    string    `json:"chat_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
