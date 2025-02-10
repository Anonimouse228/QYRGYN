package models

type Subscription struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	UserID uint    `json:"user_id" gorm:"not null"`
	Type   string  `json:"type" gorm:"not null"`
	Price  float64 `json:"price" gorm:"not null"`
	Status string  `json:"status" gorm:"not null"`
}
