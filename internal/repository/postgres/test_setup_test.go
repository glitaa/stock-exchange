package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/glitaa/stock-exchange/internal/db"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(ctx context.Context) (*sql.DB, func(), error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("stock-exchange"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, err
	}

	database, err := db.NewPostgresDB(connStr)
	if err != nil {
		return nil, nil, err
	}

	err = db.InitSchema(ctx, database)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = database.Close()
		_ = pgContainer.Terminate(context.Background())
	}

	return database, cleanup, nil
}

func clearDB(t *testing.T, database *sql.DB) {
	_, err := database.Exec(`
		TRUNCATE TABLE audit_logs RESTART IDENTITY CASCADE;
		TRUNCATE TABLE wallet_stocks RESTART IDENTITY CASCADE;
		TRUNCATE TABLE wallets RESTART IDENTITY CASCADE;
		TRUNCATE TABLE bank_stocks RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		t.Fatalf("Failed to clear db: %v", err)
	}
}
