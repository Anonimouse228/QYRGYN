package models

import "time"

type Like struct {
	ID        uint      `gorm:"primaryKey"`
	PostId    uint      `gorm:"not null"`
	UserId    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
