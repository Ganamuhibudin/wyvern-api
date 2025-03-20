package repositories

import (
	"gorm.io/gorm"
	"wyvern-api/models"
)

// TransactionRepo struct
type TransactionRepo struct {
	db *gorm.DB
}

// NewTransactionRepo initiate TransactionRepo
func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{
		db: db,
	}
}

// Insert is method to insert trx
func (repo *TransactionRepo) Insert(db *gorm.DB, transaction models.Transaction) (models.Transaction, error) {
	result := db.Create(&transaction)
	if result.Error != nil {
		return transaction, result.Error
	}

	return transaction, nil
}

func (repo *TransactionRepo) Credit(amount float64) {

}
