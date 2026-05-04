package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/glitaa/stock-exchange/internal/db"
	"github.com/glitaa/stock-exchange/internal/domain"
)

// WalletRepository is a PostgreSQL implementation of domain.WalletRepository.
type WalletRepository struct {
	dbConn *sql.DB
}

// NewWalletRepository creates a new instance of WalletRepository.
func NewWalletRepository(dbConn *sql.DB) *WalletRepository {
	return &WalletRepository{dbConn: dbConn}
}

// runner returns a transaction from context if it exists, otherwise returns the default DB connection.
func (r *WalletRepository) runner(ctx context.Context) querier {
	if tx, ok := db.GetTx(ctx); ok {
		return tx
	}
	return r.dbConn
}

// GetWallet retrieves the current state of a specific wallet.
func (r *WalletRepository) GetWallet(ctx context.Context, id string) (domain.Wallet, error) {
	var walletID string
	err := r.runner(ctx).QueryRowContext(ctx, `SELECT id FROM wallets WHERE id = $1`, id).Scan(&walletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Wallet{}, domain.ErrWalletNotFound
		}
		return domain.Wallet{}, err
	}

	query := `SELECT stock_name, quantity FROM wallet_stocks WHERE wallet_id = $1`
	rows, err := r.runner(ctx).QueryContext(ctx, query, id)
	if err != nil {
		return domain.Wallet{}, err
	}
	defer rows.Close()

	var stocks []domain.Stock
	for rows.Next() {
		var s domain.Stock
		if err := rows.Scan(&s.Name, &s.Quantity); err != nil {
			return domain.Wallet{}, err
		}
		stocks = append(stocks, s)
	}

	if err := rows.Err(); err != nil {
		return domain.Wallet{}, err
	}

	return domain.Wallet{
		ID:     walletID,
		Stocks: stocks,
	}, nil
}

// CreateWallet creates a new empty wallet.
func (r *WalletRepository) CreateWallet(ctx context.Context, id string) error {
	query := `INSERT INTO wallets (id) VALUES ($1)`
	_, err := r.runner(ctx).ExecContext(ctx, query, id)
	return err
}

// GetStockQuantity retrieves the quantity of a specific stock in a specific wallet.
func (r *WalletRepository) GetStockQuantity(ctx context.Context, walletID, stockName string) (int, error) {
	query := `SELECT quantity FROM wallet_stocks WHERE wallet_id = $1 AND stock_name = $2`

	var quantity int
	err := r.runner(ctx).QueryRowContext(ctx, query, walletID, stockName).Scan(&quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrStockNotFound
		}
		return 0, err
	}

	return quantity, nil
}

// UpdateStockQuantity modifies the quantity of a specific stock in a wallet using an UPSERT operation.
func (r *WalletRepository) UpdateStockQuantity(ctx context.Context, walletID, stockName string, delta int) error {
	query := `
		INSERT INTO wallet_stocks (wallet_id, stock_name, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (wallet_id, stock_name)
		DO UPDATE SET quantity = wallet_stocks.quantity + EXCLUDED.quantity
	`

	_, err := r.runner(ctx).ExecContext(ctx, query, walletID, stockName, delta)
	return err
}
