package main

import (
	. "ethereye/favorites"
	. "ethereye/transactions"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Main function
func main() {
	storage := NewAddressStorage("addresses.json")
	if err := storage.Load(); err != nil {
		log.Fatalf("Failed to load addresses: %v", err)
	}

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Can't read .env: %v", err)
	}
	apiKey := os.Getenv("ETHERSCAN_APT_KEY")

	http.HandleFunc("/api/v1/favorites", FavoriteAddressHandler(storage))
	http.HandleFunc("/api/v1/transactions", TransactionsHandler)
	http.HandleFunc("/api/v1/transaction-details", TransactionDetailsHandler(apiKey))
	http.HandleFunc("/api/v1/transaction-status", TransactionStatusHandler(apiKey))
	http.HandleFunc("/filtered-transactions", FilteredTransactionsHandler(apiKey))

	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
