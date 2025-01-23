package models

import "time"

type User struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Username          string    `json:"username" gorm:"not null;unique;size:20"`
	Password          string    `json:"password" gorm:"not null;size:100"`
	Email             string    `json:"email" gorm:"not null;unique"`
	Verified          bool      `gorm:"default:false"`
	VerificationToken string    `gorm:"type:varchar(255)"`
	Role              string    `json:"role" gorm:"default:'user';not null;check:role IN ('user', 'admin')"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Posts    []Post    `json:"posts" gorm:"foreignKey:UserId"`
	Comments []Comment `json:"comments" gorm:"foreignKey:UserID"`
}
