package models

import "time"

type Comment struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id" gorm:"not null;index"`
	PostID    int       `json:"post_id" gorm:"not null;index"`
	Content   string    `json:"content" gorm:"type:text;not null;size:228"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Post Post `json:"post" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}
