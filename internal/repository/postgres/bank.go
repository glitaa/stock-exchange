package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/glitaa/stock-exchange/internal/db"
	"github.com/glitaa/stock-exchange/internal/domain"
)

// querier defines the methods shared by *sql.DB and *sql.Tx.
type querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// BankRepository is a PostgreSQL implementation of domain.BankRepository.
type BankRepository struct {
	dbConn *sql.DB
}

// NewBankRepository creates a new instance of BankRepository.
func NewBankRepository(dbConn *sql.DB) *BankRepository {
	return &BankRepository{dbConn: dbConn}
}

// runner returns a transaction from context if it exists, otherwise returns the default DB connection.
func (r *BankRepository) runner(ctx context.Context) querier {
	if tx, ok := db.GetTx(ctx); ok {
		return tx
	}
	return r.dbConn
}

// GetStocks retrieves all stocks currently available in the bank.
func (r *BankRepository) GetStocks(ctx context.Context) ([]domain.Stock, error) {
	query := `SELECT name, quantity FROM bank_stocks`

	rows, err := r.runner(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []domain.Stock
	for rows.Next() {
		var s domain.Stock
		if err := rows.Scan(&s.Name, &s.Quantity); err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
}

// SetStocks replaces the entire inventory of the bank with the provided stocks.
func (r *BankRepository) SetStocks(ctx context.Context, stocks []domain.Stock) error {
	_, err := r.runner(ctx).ExecContext(ctx, `DELETE FROM bank_stocks`)
	if err != nil {
		return err
	}

	if len(stocks) == 0 {
		return nil
	}

	query := `INSERT INTO bank_stocks (name, quantity) VALUES ($1, $2)`
	for _, s := range stocks {
		_, err := r.runner(ctx).ExecContext(ctx, query, s.Name, s.Quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetStockQuantity retrieves the quantity of a specific stock in the bank.
func (r *BankRepository) GetStockQuantity(ctx context.Context, name string) (int, error) {
	query := `SELECT quantity FROM bank_stocks WHERE name = $1`

	var quantity int
	err := r.runner(ctx).QueryRowContext(ctx, query, name).Scan(&quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrStockNotFound
		}
		return 0, err
	}

	return quantity, nil
}

// UpdateStock modifies the quantity of a specific stock by a given delta.
func (r *BankRepository) UpdateStock(ctx context.Context, name string, delta int) error {
	query := `UPDATE bank_stocks SET quantity = quantity + $1 WHERE name = $2`

	result, err := r.runner(ctx).ExecContext(ctx, query, delta, name)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrStockNotFound
	}

	return nil
}
