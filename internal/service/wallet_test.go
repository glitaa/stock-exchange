package service

import (
	"context"
	"testing"

	"github.com/glitaa/stock-exchange/internal/domain"
	"github.com/glitaa/stock-exchange/internal/domain/mocks"
	"go.uber.org/mock/gomock"
)

func TestWalletService_GetWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWalletRepository(ctrl)
	svc := NewWalletService(mockRepo)
	ctx := context.Background()

	expectedWallet := domain.Wallet{
		ID:     "wallet_123",
		Stocks: []domain.Stock{{Name: "AAPL", Quantity: 10}},
	}

	mockRepo.EXPECT().GetWallet(ctx, "wallet_123").Return(expectedWallet, nil)

	w, err := svc.GetWallet(ctx, "wallet_123")
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}
	if w.ID != "wallet_123" {
		t.Errorf("expected wallet_123, got %v", w.ID)
	}
}

func TestWalletService_GetStockQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWalletRepository(ctrl)
	svc := NewWalletService(mockRepo)
	ctx := context.Background()

	t.Run("wallet exists, returns stock", func(t *testing.T) {
		mockRepo.EXPECT().GetWallet(ctx, "wallet_123").Return(domain.Wallet{}, nil)
		mockRepo.EXPECT().GetStockQuantity(ctx, "wallet_123", "AAPL").Return(10, nil)

		qty, err := svc.GetStockQuantity(ctx, "wallet_123", "AAPL")
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
		if qty != 10 {
			t.Errorf("expected 10, got %d", qty)
		}
	})

	t.Run("wallet not found", func(t *testing.T) {
		mockRepo.EXPECT().GetWallet(ctx, "wallet_123").Return(domain.Wallet{}, domain.ErrWalletNotFound)

		_, err := svc.GetStockQuantity(ctx, "wallet_123", "AAPL")
		if err != domain.ErrWalletNotFound {
			t.Fatalf("expected ErrWalletNotFound, got %v", err)
		}
	})
}
