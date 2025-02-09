package models

import "time"

type Chat struct {
	ID        int `gorm:"primaryKey"`
	UserID    int `gorm:"not null"`
	AdminID   int
	Status    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:timestamp;not null"`
	ClosedAt  time.Time `gorm:"type:timestamp"`
}
