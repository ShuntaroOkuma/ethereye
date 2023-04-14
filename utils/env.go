package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetApiKey() string {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Printf("Can't read .env: %v", err)
	}
	apiKey := os.Getenv("ETHERSCAN_APT_KEY")
	return apiKey
}
