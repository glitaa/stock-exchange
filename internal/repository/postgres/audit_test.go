package postgres

import (
	"context"
	"testing"

	"github.com/glitaa/stock-exchange/internal/domain"
)

func TestAuditLogRepository(t *testing.T) {
	ctx := context.Background()
	database, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewAuditLogRepository(database)

	t.Run("Add and GetAll", func(t *testing.T) {
		clearDB(t, database)
		
		entry := domain.LogEntry{
			Type:      domain.OperationTypeBuy,
			WalletID:  "w1",
			StockName: "AAPL",
		}
		err := repo.Add(ctx, entry)
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}

		logs, err := repo.GetAll(ctx)
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
		if len(logs) != 1 {
			t.Fatalf("expected 1 log, got %d", len(logs))
		}
		if logs[0].Type != entry.Type || logs[0].WalletID != entry.WalletID || logs[0].StockName != entry.StockName {
			t.Errorf("expected %v, got %v", entry, logs[0])
		}
	})
}
