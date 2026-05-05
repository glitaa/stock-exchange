package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/glitaa/stock-exchange/internal/db"
	"github.com/glitaa/stock-exchange/internal/handler"
	"github.com/glitaa/stock-exchange/internal/repository/postgres"
	"github.com/glitaa/stock-exchange/internal/service"
)

func main() {
	port := flag.String("port", "8080", "Port on which the server will listen")
	flag.Parse()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	dbConn, err := db.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer dbConn.Close()

	if err := db.InitSchema(context.Background(), dbConn); err != nil {
		log.Fatalf("Cannot initialize database schema: %v", err)
	}

	bankRepo := postgres.NewBankRepository(dbConn)
	walletRepo := postgres.NewWalletRepository(dbConn)
	auditRepo := postgres.NewAuditRepository(dbConn)
	txManager := db.NewTxManager(dbConn)

	bankService := service.NewBankService(bankRepo)
	walletService := service.NewWalletService(walletRepo)
	exchangeService := service.NewExchangeService(walletRepo, bankRepo, auditRepo, txManager)
	auditService := service.NewAuditService(auditRepo)

	bankHandler := handler.NewBankHandler(bankService)
	walletHandler := handler.NewWalletHandler(walletService)
	exchangeHandler := handler.NewExchangeHandler(exchangeService)
	auditHandler := handler.NewAuditHandler(auditService)
	chaosHandler := handler.NewChaosHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /wallets/{wallet_id}/stocks/{stock_name}", exchangeHandler.Trade)

	mux.HandleFunc("GET /stocks", bankHandler.GetStocks)
	mux.HandleFunc("POST /stocks", bankHandler.SetStocks)

	mux.HandleFunc("GET /wallets/{wallet_id}", walletHandler.GetWallet)
	mux.HandleFunc("GET /wallets/{wallet_id}/stocks/{stock_name}", walletHandler.GetWalletStock)

	mux.HandleFunc("GET /log", auditHandler.GetLog)

	mux.HandleFunc("POST /chaos", chaosHandler.Crash)

	log.Printf("Server is running on port %s...", *port)
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
