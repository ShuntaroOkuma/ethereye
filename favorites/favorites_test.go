package favorites

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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

func TestFavoriteAddressHandler(t *testing.T) {
	storage := NewAddressStorage("test_addresses.json")
	defer os.Remove("test_addresses.json")
	handler := FavoriteAddressHandler(storage)

	// Test adding favorite wallet address
	requestBody, _ := json.Marshal(map[string]string{"address": walletAddress})
	req := httptest.NewRequest(http.MethodPost, "/favorites?type=wallet", bytes.NewReader(requestBody))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
	}

	// Test getting favorite wallet addresses
	req = httptest.NewRequest(http.MethodGet, "/favorites?type=wallet", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	responseBody, _ := ioutil.ReadAll(rr.Body)
	var addresses []string
	json.Unmarshal(responseBody, &addresses)
	if len(addresses) != 1 || addresses[0] != walletAddress {
		t.Errorf("Unexpected favorite wallet addresses: %v", addresses)
	}
}
