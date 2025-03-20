package services

import (
	"gorm.io/gorm"
	"net/http"
	"wyvern-api/models"
	"wyvern-api/utils"
)

// TransactionProcessor interface
type TransactionProcessor interface {
	Insert(db *gorm.DB, transaction models.Transaction) (models.Transaction, error)
}

// TransactionService struct
type TransactionService struct {
	identifier int64
	tp         TransactionProcessor
	db         *gorm.DB
}

// NewTransactionService initiate TransactionService
func NewTransactionService(tp TransactionProcessor, db *gorm.DB, identifier int64) *TransactionService {
	return &TransactionService{
		identifier: identifier,
		tp:         tp,
		db:         db,
	}
}

// Credit is method for handle credit process
func (svc *TransactionService) Credit(req models.CreditRequest) models.ResponseV2 {
	log := utils.NewLoggerIdentifier("Credit", 1, svc.identifier).Service()
	log.Info("req: %+v", req)

	tx := svc.db.Begin()

	// Lock the row to prevent concurrent updates
	var user models.User
	if err := tx.Raw("SELECT * FROM users WHERE id = ? FOR UPDATE", req.UserID).Scan(&user).Error; err != nil {
		log.Warn("failed find user, error: %s", err.Error())
		errMsg := "Data not found"
		resp := models.ResponseV2{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: errMsg,
		}
		return resp
	}

	// update user balance
	if err := tx.Exec("UPDATE users SET balance = balance + ? WHERE id = ?", req.Amount, req.UserID).Error; err != nil {
		tx.Rollback()
		log.Warn("failed update user, error: %s", err.Error())
		errMsg := "Invalid Amount"
		resp := models.ResponseV2{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: errMsg,
		}
		return resp
	}

	// insert transactions
	transaction := models.Transaction{
		UserID: user.ID,
		Amount: req.Amount,
		Type:   "CREDIT",
	}
	transaction, err := svc.tp.Insert(tx, transaction)
	if err != nil {
		tx.Rollback()
		log.Warn("failed insert transaction, error: %s", err.Error())
		errMsg := "Invalid Amount"
		resp := models.ResponseV2{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: errMsg,
		}
		return resp
	}
	tx.Commit()

	resp := models.ResponseV2{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "",
		Data: models.CreditResponse{
			TransactionID: transaction.ID,
			NewBalance:    user.Balance,
		},
	}

	return resp
}

// Debit is method for handle debit process
func (svc *TransactionService) Debit(req models.DebitRequest) models.ResponseV2 {
	log := utils.NewLoggerIdentifier("Credit", 1, svc.identifier).Service()
	log.Info("req: %+v", req)

	tx := svc.db.Begin()

	// Lock the row to prevent concurrent updates
	var user models.User
	if err := tx.Raw("SELECT * FROM users WHERE id = ? FOR UPDATE", req.UserID).Scan(&user).Error; err != nil {
		log.Warn("failed find user, error: %s", err.Error())
		errMsg := "Data not found"
		resp := models.ResponseV2{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: errMsg,
		}
		return resp
	}

	// validate balance
	if user.Balance < req.Amount {
		errMsg := "Insufficient funds"
		log.Warn(errMsg)
		resp := models.ResponseV2{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: errMsg,
		}
		return resp
	}

	// update user balance
	if err := tx.Exec("UPDATE users SET balance = balance - ? WHERE id = ?", req.Amount, req.UserID).Error; err != nil {
		tx.Rollback()
		log.Warn("failed update user, error: %s", err.Error())
		errMsg := "Invalid Amount"
		resp := models.ResponseV2{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: errMsg,
		}
		return resp
	}

	// insert transactions
	transaction := models.Transaction{
		UserID: user.ID,
		Amount: req.Amount,
		Type:   "DEBIT",
	}
	transaction, err := svc.tp.Insert(tx, transaction)
	if err != nil {
		tx.Rollback()
		log.Warn("failed insert transaction, error: %s", err.Error())
		errMsg := "Invalid Amount"
		resp := models.ResponseV2{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: errMsg,
		}
		return resp
	}
	tx.Commit()

	newBalance := user.Balance - req.Amount
	resp := models.ResponseV2{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "",
		Data: models.CreditResponse{
			TransactionID: transaction.ID,
			NewBalance:    newBalance,
		},
	}

	return resp
}
