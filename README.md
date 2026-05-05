# Stock Exchange

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?logo=docker)](https://www.docker.com/)

A simplified, highly available stock exchange simulation service built in Go.

## Overview

This project simulates a basic stock market ecosystem consisting of a **Bank** and user **Wallets**. It provides a robust RESTful API to manage stock inventories, trade assets (buy/sell), and audit transactions.

The application is designed to be highly available (HA), utilizing a multi-instance Go backend load-balanced via NGINX, and backed by a PostgreSQL database with strict transactional integrity.

## Features

- **Bank:** The sole liquidity provider. Controls the available stocks.
- **Wallets:** Entities that hold stocks. Created automatically upon the first transaction.
- **Trading:** Immediate execution of buy and sell operations at a fixed price of 1.
- **Audit Log:** Tracks all successful buy and sell operations.
- **High Availability (HA):** Built-in `/chaos` endpoint to randomly crash backend instances, demonstrating Docker's auto-healing capabilities and NGINX's load balancer fault tolerance.
- **Cross-Platform:** Out-of-the-box compatibility with Windows, Linux, and macOS (across x64 and arm64 architectures).

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- Docker Compose

## Getting Started

The project provides unified startup scripts that handle building the optimized multi-stage Docker images and launching the entire cluster (PostgreSQL, NGINX, and App Replicas).

### Windows
```cmd
start.bat [PORT]
```
*(Example: `start.bat 9000`)*

### Linux / macOS
```bash
./start.sh [PORT]
```
*(Example: `./start.sh 9000`)*

If no port is provided, it defaults to `8080`. The application will be available at `http://localhost:[PORT]`.

## API Endpoints & Usage Examples

### 1. Set Bank State
Overwrites the entire inventory of the bank.
- **POST** `/stocks`
```bash
curl -X POST http://localhost:8080/stocks \
     -H "Content-Type: application/json" \
     -d '{"stocks": [{"name": "AAPL", "quantity": 1000}, {"name": "GOOG", "quantity": 500}]}'
```

### 2. Get Bank State
Returns the current state of the bank's stock inventory.
- **GET** `/stocks`

### 3. Trade Stock
Executes a buy or sell operation. If a wallet doesn't exist during a buy, it is automatically created.
- **POST** `/wallets/{wallet_id}/stocks/{stock_name}`
```bash
# Buy a stock
curl -X POST http://localhost:8080/wallets/wallet_123/stocks/AAPL \
     -H "Content-Type: application/json" \
     -d '{"type": "buy"}'

# Sell a stock
curl -X POST http://localhost:8080/wallets/wallet_123/stocks/AAPL \
     -H "Content-Type: application/json" \
     -d '{"type": "sell"}'
```

### 4. Get Wallet State
Returns the current state of a particular wallet.
- **GET** `/wallets/{wallet_id}`
```bash
curl http://localhost:8080/wallets/wallet_123
```

### 5. Get Wallet Stock Quantity
Returns the quantity of a specific stock in a specified wallet.
- **GET** `/wallets/{wallet_id}/stocks/{stock_name}`

### 6. Get Audit Log
Returns the entire audit log of successful operations in order of occurrence.
- **GET** `/log`
```bash
curl http://localhost:8080/log
```

### 7. Simulate Crash (Chaos)
Kills the instance that currently serves the request to test high availability.
- **POST** `/chaos`
```bash
curl -X POST http://localhost:8080/chaos
# Subsequent requests will succeed via NGINX routing to healthy instances.
# The Docker daemon will automatically restart the crashed container.
```
