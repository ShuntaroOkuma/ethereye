# EtherEye

EtherEye is API designed to enrich the user experience of Etherscan by providing a more feature-rich. With EtherEye, users can easily track their transactions, optimize gas fees, and manage their asset portfolios.

# Environment

- API
  - Go: '1.18.4'
- API Document
  - swagger: '1.0.0'
  - npm: '8.1.0'
  - node: '16.13.0'

# Usage

To use EtherEye, you will need to obtain an API key from Etherscan and configure it in the application. The API provides endpoints for each feature, allowing you to easily access the information you need by making HTTP requests to the appropriate URLs.

## Run API

- Setup Etherscan API Key

Sign up for an Etherscan account and get API key.

To create a new file named `.env`, run below command in root directory after you replace `YOUR_API_KEY` to your API Key.

```sh
echo 'ETHERSCAN_APT_KEY="YOUR_API_KEY"' > .env
```

- Run go

Run below command in root directory.

```sh
go run main.go
```

Your terminal will shows "Starting server on port 8080...", and now you can access `http://localhost:8080/`.

- Run API

Replace `YOUR_WALLET_ADDRESS` to you wallet adress in the following URL, and access it with you browser, you can obtain information about the transaction.

```
http://localhost:8080/api/v1/transactions?address=YOUR_WALLET_ADDRESS
```

## See API docs

Refer to the provided API documentation (Swagger UI format) for detailed information on each endpoint, including the required input parameters, expected output, and any applicable error messages.

Run below command.

```sh
cd swagger
npm install
npm start
```

Your terminal will show below:

```
> swagger@1.0.0 start
> node index.js

Swagger UI is available at http://localhost:3001
```

And access http://localhost:3001 with your browser, and you can see API docs.

# Features

## Transaction Tracker (Implemented)

1. View recent transactions related to your wallet address (1-1): Keep track of recent transactions involving your wallet address.
2. View detailed information about a specific transaction (1-2): Input a transaction ID to view detailed information such as sender, recipient, amount, gas fee, etc.
3. View real-time transaction status (1-3): Monitor the real-time status of transactions (unconfirmed, confirmed, failed).
4. Favorite specific wallet or token contract addresses for easy access (1-4): Easily access your favorite wallet addresses or token contract addresses.
5. Filter transaction history based on specific timeframes or token types (1-5): Customize your transaction history view by filtering transactions based on timeframes or token types.

## Gas Fee Optimization Tool (To be implemented)

1. View current gas fees in real-time (2-1): Stay updated on average, low, and high gas fees.
2. Receive optimal gas price suggestions (2-2): Get gas price recommendations based on the urgency of your transaction.
3. Know approximate transaction processing time (2-3): Estimate the processing time for your transaction based on the chosen gas price.
4. View past gas fee fluctuations in a graph (2-4): Understand gas fee trends by analyzing historical data.
5. See future gas fee predictions (2-5): Make informed transaction timing decisions based on future gas fee predictions.

## Asset Portfolio Manager (To be implemented)

1. View current price and total value of tokens in a wallet (3-1): Input a wallet address with multiple tokens and view their current price and total value.
2. View portfolio performance in a time-series graph (3-2): Monitor changes in the value of your portfolio over time.
3. View trade history of a specific token (3-3): Compare the acquisition price and the current price of a specific token.
4. View asset allocation in a pie chart (3-4): See the proportion of each token in your portfolio.
5. Manage assets across multiple wallets (3-5): Add different wallet addresses and manage the asset situation in a unified manner.
