package postgres

import (
	"context"
	"testing"

	"github.com/glitaa/stock-exchange/internal/domain"
)

func TestWalletRepository(t *testing.T) {
	ctx := context.Background()
	database, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewWalletRepository(database)

	t.Run("Create and Get Wallet", func(t *testing.T) {
		clearDB(t, database)

		err := repo.CreateWallet(ctx, "w1")
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}

		w, err := repo.GetWallet(ctx, "w1")
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
		if w.ID != "w1" {
			t.Errorf("expected w1, got %s", w.ID)
		}
	})

	t.Run("Get missing wallet", func(t *testing.T) {
		clearDB(t, database)
		_, err := repo.GetWallet(ctx, "missing")
		if err != domain.ErrWalletNotFound {
			t.Fatalf("expected ErrWalletNotFound, got %v", err)
		}
	})

	t.Run("Update and Get Stock Quantity", func(t *testing.T) {
		clearDB(t, database)
		if err := repo.CreateWallet(ctx, "w1"); err != nil {
			t.Fatalf("failed to create wallet: %v", err)
		}

		err := repo.UpdateStockQuantity(ctx, "w1", "AAPL", 10)
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}

		qty, err := repo.GetStockQuantity(ctx, "w1", "AAPL")
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
		if qty != 10 {
			t.Errorf("expected 10, got %d", qty)
		}

		// Update again
		err = repo.UpdateStockQuantity(ctx, "w1", "AAPL", -5)
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}

		qty, err = repo.GetStockQuantity(ctx, "w1", "AAPL")
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
		if qty != 5 {
			t.Errorf("expected 5, got %d", qty)
		}
	})
}
