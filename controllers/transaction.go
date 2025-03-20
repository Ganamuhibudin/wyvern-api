package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"wyvern-api/config"
	"wyvern-api/models"
	"wyvern-api/repositories"
	"wyvern-api/services"
	"wyvern-api/utils"
)

type TransactionController struct {
}

// NewTransactionController initiate TransactionController
func NewTransactionController() *TransactionController {
	return &TransactionController{}
}

// Credit is method to increase user balance
func (c *TransactionController) Credit(ctx *gin.Context) {
	identifier := time.Now().UnixNano()
	log := utils.NewLoggerIdentifier("Credit", 0, identifier).Controller().Start()

	var req models.CreditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn("bad request, error: %s", err.Error())
		log.End()
		resp := models.ResponseV2{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid amount",
		}
		ctx.JSON(resp.Code, resp)
		return
	}

	db := config.DB
	transactionRepo := repositories.NewTransactionRepo(db)
	svc := services.NewTransactionService(transactionRepo, db, identifier)
	response := svc.Credit(req)

	log.End()
	ctx.JSON(http.StatusOK, response)
}

// Debit is method to deduct user balance
func (c *TransactionController) Debit(ctx *gin.Context) {
	identifier := time.Now().UnixNano()
	log := utils.NewLoggerIdentifier("Debit", 0, identifier).Controller().Start()

	var req models.DebitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn("bad request, error: %s", err.Error())
		log.End()
		resp := models.ResponseV2{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid amount",
		}
		ctx.JSON(resp.Code, resp)
		return
	}

	db := config.DB
	transactionRepo := repositories.NewTransactionRepo(db)
	svc := services.NewTransactionService(transactionRepo, db, identifier)
	response := svc.Debit(req)

	log.End()
	ctx.JSON(http.StatusOK, response)
}
