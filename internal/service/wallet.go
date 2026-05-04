package service

import (
	"context"

	"github.com/glitaa/stock-exchange/internal/domain"
)

// WalletService handles read-only business logic for wallets.
type WalletService struct {
	walletRepo domain.WalletRepository
}

// NewWalletService creates a new instance of WalletService.
func NewWalletService(walletRepo domain.WalletRepository) *WalletService {
	return &WalletService{walletRepo: walletRepo}
}

// GetWallet retrieves wallet information by ID.
func (s *WalletService) GetWallet(ctx context.Context, walletID string) (domain.Wallet, error) {
	return s.walletRepo.GetWallet(ctx, walletID)
}

func (s *WalletService) GetStockQuantity(ctx context.Context, walletID, stockName string) (int, error) {
	_, err := s.walletRepo.GetWallet(ctx, walletID)
	if err != nil {
		return 0, err
	}
	return s.walletRepo.GetStockQuantity(ctx, walletID, stockName)
}
