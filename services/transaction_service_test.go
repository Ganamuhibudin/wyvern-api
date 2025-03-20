package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
	"wyvern-api/models"
)

var DBMock *gorm.DB

const (
	DB_USERNAME = "root"
	DB_PASSWORD = ""
	DB_URL      = "127.0.0.1"
	DB_PORT     = "3306"
	DB_DATABASE = "wyvern-api"
	DB_DEBUG    = true
)

func MockLoadDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USERNAME, DB_PASSWORD, DB_URL, DB_PORT, DB_DATABASE)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DBMock = db
}

func TestTransactionService_Concurrent(t *testing.T) {
	MockLoadDatabase()
	user := models.User{Username: "Fulan"}
	DBMock.Create(&user)

	apiURL := "http://127.0.0.1:8080/api/transactions/credit"
	var wg sync.WaitGroup
	client := &http.Client{Timeout: 10 * time.Second}

	creditAmount := 1000
	numRequests := 100     // Number of concurrent requests
	numTransactions := 100 // each request do 100 trx
	wg.Add(numRequests)    // Add all Goroutines to WaitGroup

	req := models.CreditRequest{
		UserID: user.ID,
		Amount: float64(creditAmount),
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numTransactions; j++ {
				// Create a new request
				req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
				if err != nil {
					t.Errorf("Request creation failed: %v", err)
					return
				}

				// Send request
				resp, err := client.Do(req)
				if err != nil {
					t.Errorf("Request %d failed: %v", i, err)
					return
				}
				defer resp.Body.Close()

				// Validate response status code
				if resp.StatusCode != http.StatusOK {
					t.Errorf("Request %d failed: expected 200, got %d", i, resp.StatusCode)
				}
			}
		}()
	}

	wg.Wait()

	// validate balance
	var updatedUser models.User
	DBMock.First(&updatedUser, user.ID)

	expectedBalance := float64(numRequests * numTransactions * creditAmount) // 100 * 100 * 1000 = 10.000.000
	if updatedUser.Balance != expectedBalance {
		t.Errorf("Balance does not match: expected %v, got %v", expectedBalance, updatedUser.Balance)
	}
	fmt.Printf("User balance is match!")
}
