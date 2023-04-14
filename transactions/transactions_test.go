package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

var walletAddress string
var transactionID string

func LoadTestEnv() {
	err := godotenv.Load("../.env.test")
	if err != nil {
		fmt.Printf("Can't read .env.test: %v", err)
	}
	walletAddress = os.Getenv("WALLET_ADDRESS")
	transactionID = os.Getenv("TRANSACTION_ID")
}

/************
common
************/
func serveHTTPTransactionsHandler(method string, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TransactionsHandler)
	handler.ServeHTTP(rr, req)

	return rr
}

/************
test body
************/
func TestFetchTransactions(t *testing.T) {
	LoadEnv()
	LoadTestEnv()

	startIndex := 1
	count := 3
	transactions, err := FetchTransactions(walletAddress, startIndex, count, apiKey)

	if err != nil {
		t.Errorf("fetchTransactions() returned an error: %v", err)
	}

	if len(transactions) != 3 {
		t.Error("fetchTransactions() returned an empty result")
	}

}

func TestTransactionsHandler(t *testing.T) {
	LoadTestEnv()
	t.Run("Test case 1: Valid wallet address", func(t *testing.T) {
		rr := serveHTTPTransactionsHandler("GET", "/api/v1/transactions?address="+walletAddress)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("TransactionsHandler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		var transactions []Transaction
		err := json.NewDecoder(rr.Body).Decode(&transactions)

		if err != nil {
			t.Errorf("Failed to decode response JSON: %v", err)
		}

		if len(transactions) == 0 {
			t.Error("TransactionsHandler returned an empty result")
		}
	})

	t.Run("Test case 2: Invalid wallet address", func(t *testing.T) {
		rr := serveHTTPTransactionsHandler("GET", "/transactions?address=invalid_address")

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
		}
	})

	t.Run("Test case 3: Test with missing wallet address", func(t *testing.T) {
		rr := serveHTTPTransactionsHandler("GET", "/api/v1/transactions")

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("TransactionsHandler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}

		if !strings.Contains(rr.Body.String(), "Missing 'address' query parameter") {
			t.Error("TransactionsHandler did not return expected error message for missing address")
		}
	})

	t.Run("Test case 4: Transactions filtering with startIndex and count", func(t *testing.T) {
		rr := serveHTTPTransactionsHandler("GET", "/api/v1/transactions?address="+walletAddress+"&startIndex=1&count=3")

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("TransactionsHandler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var transactions []Transaction
		err := json.NewDecoder(rr.Body).Decode(&transactions)

		if err != nil {
			t.Errorf("Failed to decode response JSON: %v", err)
		}

		if len(transactions) != 3 {
			t.Errorf("Expected 3 transactions, got %d", len(transactions))
		}
	})

}

func TestFetchTransactionDetails(t *testing.T) {
	LoadEnv()
	LoadTestEnv()

	transactionDetails, err := FetchTransactionDetails(transactionID, apiKey)

	if err != nil {
		t.Fatalf("FetchTransactionDetails failed: %v", err)
	}

	if transactionDetails.From == "" || transactionDetails.To == "" || transactionDetails.Value == "" || transactionDetails.Gas == "" || transactionDetails.GasPrice == "" {
		t.Errorf("transactionDetails has empty fields")
	}
}

func TestTransactionDetailsHandler(t *testing.T) {
	LoadTestEnv()

	handler := TransactionDetailsHandler(apiKey)

	// Test case 1: Valid transaction ID
	req := httptest.NewRequest(http.MethodGet, "/transaction_details?txid="+transactionID, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("transactionDetailsHandler returned non-200 status code: %d", rec.Code)
	}

	var transactionDetails TransactionDetails
	err := json.Unmarshal(rec.Body.Bytes(), &transactionDetails)
	if err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if transactionDetails.From == "" || transactionDetails.To == "" || transactionDetails.Value == "" || transactionDetails.Gas == "" || transactionDetails.GasPrice == "" {
		t.Errorf("transactionDetails has empty fields")
	}

	// Test case 2: Missing transaction ID
	req = httptest.NewRequest(http.MethodGet, "/transaction_details", nil)
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("transactionDetailsHandler returned non-400 status code: %d", rec.Code)
	}
}

func TestFetchTransactionStatus(t *testing.T) {
	LoadEnv()
	LoadTestEnv()

	t.Run("Test with a valid transaction ID", func(t *testing.T) {
		txStatus, err := FetchTransactionStatus(transactionID, apiKey)
		if err != nil {
			t.Errorf("Failed to fetch transaction status for valid transaction ID: %v", err)
			return
		}
		if txStatus.Status == "" {
			t.Errorf("Transaction status is empty for valid transaction ID")
		}

		if txStatus.Status != "1" {
			t.Errorf("Expected transaction status is 1, but got %v", txStatus.Status)
		}
	})

	t.Run("Test with an invalid transaction ID", func(t *testing.T) {
		invalidTxID := "0xINVALID_TRANSACTION_ID"
		txStatus, _ := FetchTransactionStatus(invalidTxID, apiKey)
		if txStatus.Status != "" {
			t.Errorf("Expected status is vacant, but got %v", txStatus.Status)
		}
	})

	t.Run("Test with an empty transaction ID", func(t *testing.T) {
		emptyTxID := ""
		_, err := FetchTransactionStatus(emptyTxID, apiKey)
		if err == nil {
			t.Errorf("Expected error for empty transaction ID, but got none")
		}
	})

	t.Run("Test with an invalid API key", func(t *testing.T) {
		invalidAPIKey := "INVALID_API_KEY"
		_, err := FetchTransactionStatus(transactionID, invalidAPIKey)
		if err == nil {
			t.Errorf("Expected error for invalid API key, but got none")
		}
	})
}

func TestFetchFilteredTransactions(t *testing.T) {
	LoadEnv()
	LoadTestEnv()

	startDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC)
	tokenType := "ETH" // Example token type

	transactions, err := FetchFilteredTransactions(walletAddress, apiKey, &startDate, &endDate, tokenType)
	if err != nil {
		t.Errorf("Error fetching transactions: %v", err)
	}

	for _, tx := range transactions {
		if tx.Timestamp.Before(startDate) || tx.Timestamp.After(endDate) {
			t.Errorf("Transaction is not within specified date range")
		}

		if tokenType != "" && tx.TokenType != tokenType {
			t.Errorf("Transaction token type does not match specified token type")
		}
	}
}
