package models

import "time"

type Post struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	UserId     int       `json:"userid"`
	Content    string    `json:"content"`
	Likes      int       `json:"likes"`
	Created_at time.Time `json:"createdat"`
	Updated_at time.Time `json:"updatedat"`
}
