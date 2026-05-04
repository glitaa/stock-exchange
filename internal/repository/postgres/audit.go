package postgres

import (
	"context"
	"database/sql"

	"github.com/glitaa/stock-exchange/internal/db"
	"github.com/glitaa/stock-exchange/internal/domain"
)

// AuditRepository is a PostgreSQL implementation of domain.AuditLogRepository.
type AuditRepository struct {
	dbConn *sql.DB
}

// NewAuditRepository creates a new instance of AuditRepository.
func NewAuditRepository(dbConn *sql.DB) *AuditRepository {
	return &AuditRepository{dbConn: dbConn}
}

// runner returns a transaction from context if it exists, otherwise returns the default DB connection.
func (r *AuditRepository) runner(ctx context.Context) querier {
	if tx, ok := db.GetTx(ctx); ok {
		return tx
	}
	return r.dbConn
}

// Add inserts a new log entry into the audit_logs table.
func (r *AuditRepository) Add(ctx context.Context, entry domain.LogEntry) error {
	query := `
		INSERT INTO audit_logs (operation_type, wallet_id, stock_name)
		VALUES ($1, $2, $3)
	`
	_, err := r.runner(ctx).ExecContext(ctx, query, entry.Type, entry.WalletID, entry.StockName)
	return err
}

// GetAll retrieves all audit log entries, ordered by creation time.
func (r *AuditRepository) GetAll(ctx context.Context) ([]domain.LogEntry, error) {
	query := `
		SELECT operation_type, wallet_id, stock_name
		FROM audit_logs
		ORDER BY created_at ASC
	`
	rows, err := r.runner(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.LogEntry
	for rows.Next() {
		var entry domain.LogEntry
		if err := rows.Scan(&entry.Type, &entry.WalletID, &entry.StockName); err != nil {
			return nil, err
		}
		logs = append(logs, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
