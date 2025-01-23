package models

import "time"

type Post struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserId    int       `json:"user_id" gorm:"not null;index"`
	Content   string    `json:"content" gorm:"type:text;not null;size:228"`
	Likes     int       `json:"likes" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}
