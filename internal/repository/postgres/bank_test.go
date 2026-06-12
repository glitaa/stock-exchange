package postgres

import (
	"context"
	"testing"

	"github.com/glitaa/stock-exchange/internal/domain"
)

func TestBankRepository(t *testing.T) {
	ctx := context.Background()
	database, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewBankRepository(database)

	t.Run("Set and Get Stocks", func(t *testing.T) {
		clearDB(t, database)
		stocks := []domain.Stock{{Name: "AAPL", Quantity: 100}, {Name: "GOOG", Quantity: 50}}
		err := repo.SetStocks(ctx, stocks)
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}

		res, err := repo.GetStocks(ctx)
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
		if len(res) != 2 {
			t.Errorf("expected 2 stocks, got %d", len(res))
		}
	})

	t.Run("Update Stock Quantity", func(t *testing.T) {
		clearDB(t, database)
		stocks := []domain.Stock{{Name: "AAPL", Quantity: 100}}
		repo.SetStocks(ctx, stocks)

		err := repo.UpdateStockQuantity(ctx, "AAPL", -10)
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}

		qty, err := repo.GetStockQuantity(ctx, "AAPL")
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
		if qty != 90 {
			t.Errorf("expected 90, got %d", qty)
		}
	})
}
