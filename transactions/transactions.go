package transactions

import (
	"encoding/json"
	"errors"
	. "ethereye/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Transaction struct definition
type Transaction struct {
	ID          string    `json:"hash"`
	FromAddress string    `json:"from"`
	ToAddress   string    `json:"to"`
	Value       string    `json:"value"`
	GasPrice    string    `json:"gasPrice"`
	TokenType   string    `json:"tokenType"`
	BlockHeight uint64    `json:"blockHeight"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timeStamp"`
}

type TransactionDetails struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	Gas       string `json:"gas"`
	GasPrice  string `json:"gasPrice"`
	InputData string `json:"inputData"`
}

type TransactionStatus struct {
	TxID   string `json:"txid"`
	Status string `json:"status"`
}

type TokenTransfer struct {
	BlockNumber  int64   `json:"blockNumber"`
	Timestamp    int64   `json:"timeStamp"`
	Hash         string  `json:"hash"`
	From         string  `json:"from"`
	To           string  `json:"to"`
	Value        float64 `json:"value"`
	ContractAddr string  `json:"contractAddress"`
}

/******************
Transactions
******************/

func FetchTransactions(apiKey, walletAddress string) ([]Transaction, error) {
	var result []Transaction

	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s", walletAddress, apiKey)

	// Send the HTTP request to the Etherscan API
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	if jsonData["status"].(string) != "1" {
		return nil, fmt.Errorf("error fetching transactions: %s", jsonData)
	}

	transactions := jsonData["result"].([]interface{})

	for _, tx := range transactions {
		txData := tx.(map[string]interface{})
		blockHeight, _ := strconv.Atoi(txData["blockNumber"].(string))
		timestamp, _ := strconv.Atoi(txData["timeStamp"].(string))

		transaction := Transaction{
			ID:          txData["hash"].(string),
			FromAddress: txData["from"].(string),
			ToAddress:   txData["to"].(string),
			Value:       txData["value"].(string),
			GasPrice:    txData["gasPrice"].(string),
			TokenType:   "ETH",
			BlockHeight: uint64(blockHeight),
			Status:      txData["txreceipt_status"].(string),
			Timestamp:   time.Unix(int64(timestamp), 0),
		}
		result = append(result, transaction)
	}

	return result, nil
}

// Transactions API handler
func TransactionsHandler(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the query parameters
		walletAddress := r.URL.Query().Get("address")
		if walletAddress == "" {
			http.Error(w, "Missing 'address' query parameter", http.StatusBadRequest)
			return
		}

		apiKey := GetApiKey()

		// Fetch the transactions
		transactions, err := FetchTransactions(apiKey, walletAddress)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching transactions: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// Write the JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transactions)
	}
}

/******************
Transaction Details
******************/

func FetchTransactionDetails(transactionID string, apiKey string) (TransactionDetails, error) {
	apiUrl := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getTransactionByHash&txhash=%s&apikey=%s", transactionID, apiKey)

	// Send the HTTP request to the Etherscan API
	resp, err := http.Get(apiUrl)
	if err != nil {
		return TransactionDetails{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TransactionDetails{}, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return TransactionDetails{}, err
	}

	transaction := data["result"].(map[string]interface{})
	return TransactionDetails{
		From:      transaction["from"].(string),
		To:        transaction["to"].(string),
		Value:     transaction["value"].(string),
		Gas:       transaction["gas"].(string),
		GasPrice:  transaction["gasPrice"].(string),
		InputData: transaction["input"].(string),
	}, nil
}

func TransactionDetailsHandler(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactionID := r.URL.Query().Get("txid")
		if transactionID == "" {
			http.Error(w, "missing transaction ID", http.StatusBadRequest)
			return
		}

		transactionDetails, err := FetchTransactionDetails(transactionID, apiKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(transactionDetails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

/******************
Transaction Status
******************/

func FetchTransactionStatus(txID, apiKey string) (TransactionStatus, error) {
	if txID == "" {
		return TransactionStatus{}, errors.New("empty transaction ID")
	}

	apiUrl := fmt.Sprintf("https://api.etherscan.io/api?module=transaction&action=gettxreceiptstatus&txhash=%s&apikey=%s", txID, apiKey)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return TransactionStatus{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TransactionStatus{}, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return TransactionStatus{}, err
	}

	etherscanStatus, ok := data["status"].(string)
	if !ok {
		return TransactionStatus{}, errors.New("could not parse transaction status")
	}

	if etherscanStatus == "0" {
		return TransactionStatus{}, errors.New("Error:" + data["result"].(string))
	}

	return TransactionStatus{
		TxID:   txID,
		Status: data["result"].(map[string]interface{})["status"].(string),
	}, nil
}

// HTTP handler for fetching transaction status
func TransactionStatusHandler(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txID := r.URL.Query().Get("txid")
		if txID == "" {
			http.Error(w, "missing txid parameter", http.StatusBadRequest)
			return
		}

		txStatus, err := FetchTransactionStatus(txID, apiKey)
		if err != nil {
			http.Error(w, "could not fetch transaction status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(txStatus)
	}
}

/******************
Filtering Transaction
******************/
func FetchFilteredTransactions(walletAddress string, apiKey string, startDate *time.Time, endDate *time.Time, tokenType string) ([]Transaction, error) {
	// Fetch transactions using Etherscan API
	transactions, err := FetchTransactions(apiKey, walletAddress)
	if err != nil {
		return nil, err
	}

	// If date range is specified, filter transactions by date range
	if startDate != nil && endDate != nil {
		startTimestamp := startDate.Unix()
		endTimestamp := endDate.Unix()

		filteredTransactions := make([]Transaction, 0)
		for _, tx := range transactions {
			if tx.Timestamp.Unix() >= startTimestamp && tx.Timestamp.Unix() <= endTimestamp {
				filteredTransactions = append(filteredTransactions, tx)
			}
		}
		transactions = filteredTransactions
	}

	// If tokenType is specified, filter transactions by token type
	if tokenType != "" {
		filteredByTokenType := make([]Transaction, 0)
		for _, tx := range transactions {
			if tx.TokenType == tokenType {
				filteredByTokenType = append(filteredByTokenType, tx)
			}
		}
		transactions = filteredByTokenType
	}

	return transactions, nil
}

type FilteredTransactionsRequest struct {
	WalletAddress string `json:"wallet_address"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	TokenType     string `json:"token_type"`
}

func FilteredTransactionsHandler(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request FilteredTransactionsRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Failed to parse request", http.StatusBadRequest)
			return
		}

		var startDate, endDate *time.Time

		// Convert date strings to time.Time pointers
		if request.StartDate != "" {
			parsedStartDate, err := time.Parse(time.RFC3339, request.StartDate)
			if err != nil {
				http.Error(w, "Invalid start date format", http.StatusBadRequest)
				return
			}
			startDate = &parsedStartDate
		}

		if request.EndDate != "" {
			parsedEndDate, err := time.Parse(time.RFC3339, request.EndDate)
			if err != nil {
				http.Error(w, "Invalid end date format", http.StatusBadRequest)
				return
			}
			endDate = &parsedEndDate
		}

		transactions, err := FetchFilteredTransactions(request.WalletAddress, apiKey, startDate, endDate, request.TokenType)
		if err != nil {
			http.Error(w, "Failed to fetch filtered transactions", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(transactions)
		if err != nil {
			http.Error(w, "Failed to encode transactions", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
