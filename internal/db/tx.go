package db

import (
	"context"
	"database/sql"
	"fmt"
)

type contextKey string

const txKey contextKey = "tx"

// TxManager implements the domain.TxManager interface for a relational database (PostgreSQL).
type TxManager struct {
	db *sql.DB
}

// NewTxManager creates a new instance of TxManager.
func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

// RunInTx executes the provided function within a single database transaction.
func (m *TxManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Inject the transaction object into the context for repositories to use.
	txCtx := context.WithValue(ctx, txKey, tx)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTx retrieves the active *sql.Tx from the context (if it exists).
// Repositories will use this function to check if they are running inside a transaction.
func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	return tx, ok
}
