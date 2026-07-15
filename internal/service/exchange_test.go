package service

import (
	"context"
	"testing"

	"github.com/glitaa/stock-exchange/internal/domain"
	"github.com/glitaa/stock-exchange/internal/domain/mocks"
	"go.uber.org/mock/gomock"
)

func TestExchangeService_BuyStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mocks.NewMockWalletRepository(ctrl)
	mockBankRepo := mocks.NewMockBankRepository(ctrl)
	mockAuditRepo := mocks.NewMockAuditLogRepository(ctrl)
	mockTxManager := mocks.NewMockTxManager(ctrl)

	svc := NewExchangeService(mockWalletRepo, mockBankRepo, mockAuditRepo, mockTxManager)
	ctx := context.Background()

	t.Run("success - wallet exists", func(t *testing.T) {
		mockWalletRepo.EXPECT().GetWallet(ctx, "w1").Return(domain.Wallet{}, nil)

		// Mock TxManager to just execute the function
		mockTxManager.EXPECT().RunInTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

		mockBankRepo.EXPECT().GetStockQuantity(ctx, "AAPL").Return(10, nil)
		mockBankRepo.EXPECT().UpdateStockQuantity(ctx, "AAPL", -1).Return(nil)
		mockWalletRepo.EXPECT().UpdateStockQuantity(ctx, "w1", "AAPL", 1).Return(nil)
		mockAuditRepo.EXPECT().Add(ctx, domain.LogEntry{Type: domain.OperationTypeBuy, WalletID: "w1", StockName: "AAPL"}).Return(nil)

		err := svc.BuyStock(ctx, "w1", "AAPL")
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
	})

	t.Run("success - wallet created", func(t *testing.T) {
		mockWalletRepo.EXPECT().GetWallet(ctx, "w2").Return(domain.Wallet{}, domain.ErrWalletNotFound)
		mockWalletRepo.EXPECT().CreateWallet(ctx, "w2").Return(nil)

		mockTxManager.EXPECT().RunInTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

		mockBankRepo.EXPECT().GetStockQuantity(ctx, "AAPL").Return(5, nil)
		mockBankRepo.EXPECT().UpdateStockQuantity(ctx, "AAPL", -1).Return(nil)
		mockWalletRepo.EXPECT().UpdateStockQuantity(ctx, "w2", "AAPL", 1).Return(nil)
		mockAuditRepo.EXPECT().Add(ctx, domain.LogEntry{Type: domain.OperationTypeBuy, WalletID: "w2", StockName: "AAPL"}).Return(nil)

		err := svc.BuyStock(ctx, "w2", "AAPL")
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
	})

	t.Run("failure - bank stock empty", func(t *testing.T) {
		mockWalletRepo.EXPECT().GetWallet(ctx, "w1").Return(domain.Wallet{}, nil)

		mockTxManager.EXPECT().RunInTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

		mockBankRepo.EXPECT().GetStockQuantity(ctx, "AAPL").Return(0, nil)

		err := svc.BuyStock(ctx, "w1", "AAPL")
		if err != domain.ErrInsufficientStock {
			t.Fatalf("expected ErrInsufficientStock, got %v", err)
		}
	})

	t.Run("failure - stock not found in bank", func(t *testing.T) {
		mockWalletRepo.EXPECT().GetWallet(ctx, "w1").Return(domain.Wallet{}, nil)

		mockTxManager.EXPECT().RunInTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

		mockBankRepo.EXPECT().GetStockQuantity(ctx, "NOPE").Return(0, domain.ErrStockNotFound)

		err := svc.BuyStock(ctx, "w1", "NOPE")
		if err != domain.ErrStockNotFound {
			t.Fatalf("expected ErrStockNotFound, got %v", err)
		}
	})
}

func TestExchangeService_SellStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mocks.NewMockWalletRepository(ctrl)
	mockBankRepo := mocks.NewMockBankRepository(ctrl)
	mockAuditRepo := mocks.NewMockAuditLogRepository(ctrl)
	mockTxManager := mocks.NewMockTxManager(ctrl)

	svc := NewExchangeService(mockWalletRepo, mockBankRepo, mockAuditRepo, mockTxManager)
	ctx := context.Background()

	t.Run("success - sell stock", func(t *testing.T) {
		mockWalletRepo.EXPECT().GetWallet(ctx, "w1").Return(domain.Wallet{}, nil)

		mockTxManager.EXPECT().RunInTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

		mockBankRepo.EXPECT().GetStockQuantity(ctx, "AAPL").Return(10, nil)
		mockWalletRepo.EXPECT().GetStockQuantity(ctx, "w1", "AAPL").Return(1, nil)
		mockWalletRepo.EXPECT().UpdateStockQuantity(ctx, "w1", "AAPL", -1).Return(nil)
		mockBankRepo.EXPECT().UpdateStockQuantity(ctx, "AAPL", 1).Return(nil)
		mockAuditRepo.EXPECT().Add(ctx, domain.LogEntry{Type: domain.OperationTypeSell, WalletID: "w1", StockName: "AAPL"}).Return(nil)

		err := svc.SellStock(ctx, "w1", "AAPL")
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
	})

	t.Run("failure - wallet stock empty", func(t *testing.T) {
		mockWalletRepo.EXPECT().GetWallet(ctx, "w1").Return(domain.Wallet{}, nil)

		mockTxManager.EXPECT().RunInTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

		mockBankRepo.EXPECT().GetStockQuantity(ctx, "AAPL").Return(10, nil)
		mockWalletRepo.EXPECT().GetStockQuantity(ctx, "w1", "AAPL").Return(0, nil)

		err := svc.SellStock(ctx, "w1", "AAPL")
		if err != domain.ErrInsufficientStock {
			t.Fatalf("expected ErrInsufficientStock, got %v", err)
		}
	})

	t.Run("failure - stock not found in wallet", func(t *testing.T) {
		mockWalletRepo.EXPECT().GetWallet(ctx, "w1").Return(domain.Wallet{}, nil)

		mockTxManager.EXPECT().RunInTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

		mockBankRepo.EXPECT().GetStockQuantity(ctx, "AAPL").Return(10, nil)
		mockWalletRepo.EXPECT().GetStockQuantity(ctx, "w1", "AAPL").Return(0, domain.ErrStockNotFound)

		err := svc.SellStock(ctx, "w1", "AAPL")
		if err != domain.ErrInsufficientStock {
			t.Fatalf("expected ErrInsufficientStock, got %v", err)
		}
	})
}
