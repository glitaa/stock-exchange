package service

import (
	"context"

	"github.com/glitaa/stock-exchange/internal/domain"
)

// WalletService handles read-only business logic for wallets.
type WalletService struct {
	repo domain.WalletRepository
}

// NewWalletService creates a new instance of WalletService.
func NewWalletService(repo domain.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

// GetWallet retrieves wallet information by ID.
func (s *WalletService) GetWallet(ctx context.Context, id string) (domain.Wallet, error) {
	return s.repo.GetWallet(ctx, id)
}

// GetStockQuantity retrieves the quantity of a specific stock in a wallet.
func (s *WalletService) GetStockQuantity(ctx context.Context, walletID, stockName string) (int, error) {
	_, err := s.repo.GetWallet(ctx, walletID)
	if err != nil {
		return 0, err
	}
	return s.repo.GetStockQuantity(ctx, walletID, stockName)
}
