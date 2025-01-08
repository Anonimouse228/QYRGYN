package models

import "time"

type User struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Username          string    `json:"username"`
	Password          string    `json:"password"`
	Email             string    `json:"email"`
	Verified          bool      `gorm:"default:false"`            // For email verification
	VerificationToken string    `gorm:"type:varchar(255);unique"` // Token
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
