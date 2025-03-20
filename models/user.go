package models

import "time"

// User struct user
type User struct {
	ID        int64     `gorm:"column:id" json:"id"`
	Username  string    `gorm:"column:username" json:"username"`
	Balance   float64   `gorm:"column:balance" json:"balance"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
