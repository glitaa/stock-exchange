package service

import (
	"context"

	"github.com/glitaa/stock-exchange/internal/domain"
)

// BankService handles business logic related to the bank's stock inventory.
type BankService struct {
	repo domain.BankRepository
}

// NewBankService creates a new instance of BankService.
func NewBankService(repo domain.BankRepository) *BankService {
	return &BankService{repo: repo}
}

// GetStocks retrieves the current stock inventory from the bank.
func (s *BankService) GetStocks(ctx context.Context) ([]domain.Stock, error) {
	return s.repo.GetStocks(ctx)
}

func (s *BankService) SetStocks(ctx context.Context, stocks []domain.Stock) error {
	return s.repo.SetStocks(ctx, stocks)
}
