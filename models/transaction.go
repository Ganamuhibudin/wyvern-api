package models

import "time"

// Transaction struct transaction table
type Transaction struct {
	ID        int64     `gorm:"column:id" json:"id"`
	UserID    int64     `gorm:"column:user_id" json:"user_id"`
	Amount    float64   `gorm:"column:amount" json:"amount"`
	Type      string    `gorm:"column:type" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// CreditRequest struct for credit request
type CreditRequest struct {
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

// CreditResponse for credit response
type CreditResponse struct {
	TransactionID int64   `json:"transaction_id"`
	NewBalance    float64 `json:"new_balance"`
}

// DebitRequest struct for debit request
type DebitRequest struct {
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

// DebitResponse for debit response
type DebitResponse struct {
	TransactionID int64   `json:"transaction_id"`
	NewBalance    float64 `json:"new_balance"`
}
