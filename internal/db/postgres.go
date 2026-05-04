package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgresDB establishes a connection to the PostgreSQL database.
func NewPostgresDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL database: %w", err)
	}

	return db, nil
}

// InitSchema creates the necessary tables in the database if they do not exist.
func InitSchema(ctx context.Context, db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS bank_stocks (
		name VARCHAR(255) PRIMARY KEY,
		quantity INT NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS wallets (
		id VARCHAR(255) PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS wallet_stocks (
		wallet_id VARCHAR(255) REFERENCES wallets(id) ON DELETE CASCADE,
		stock_name VARCHAR(255) NOT NULL,
		quantity INT NOT NULL DEFAULT 0,
		PRIMARY KEY (wallet_id, stock_name)
	);

	CREATE TABLE IF NOT EXISTS audit_logs (
		id SERIAL PRIMARY KEY,
		operation_type VARCHAR(50) NOT NULL,
		wallet_id VARCHAR(255) NOT NULL,
		stock_name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}

	return nil
}
