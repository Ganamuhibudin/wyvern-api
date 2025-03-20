package routers

import (
	"github.com/gin-gonic/gin"
	"wyvern-api/controllers"
)

// Routes initiate router
func Routes(route *gin.Engine) {
	transactionController := controllers.NewTransactionController()

	api := route.Group("/api")
	transaction := api.Group("/transactions")
	transaction.POST("/credit", transactionController.Credit)
	transaction.POST("/debit", transactionController.Debit)
}
