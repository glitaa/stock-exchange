package service

import (
	"context"
	"errors"

	"github.com/glitaa/stock-exchange/internal/domain"
)

// ExchangeService handles trading operations between wallets and the bank.
type ExchangeService struct {
	walletRepo domain.WalletRepository
	bankRepo   domain.BankRepository
	auditRepo  domain.AuditLogRepository
	txManager  domain.TxManager
}

// NewExchangeService creates a new instance of ExchangeService.
func NewExchangeService(walletRepo domain.WalletRepository, bankRepo domain.BankRepository, auditRepo domain.AuditLogRepository, txManager domain.TxManager) *ExchangeService {
	return &ExchangeService{
		walletRepo: walletRepo,
		bankRepo:   bankRepo,
		auditRepo:  auditRepo,
		txManager:  txManager,
	}
}

// BuyStock executes a purchase of a single stock from the bank to the wallet.
func (s *ExchangeService) BuyStock(ctx context.Context, walletID string, stockName string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		bankQty, err := s.ensureTradingPrerequisites(txCtx, walletID, stockName)
		if err != nil {
			return err
		}

		if bankQty <= 0 {
			return domain.ErrInsufficientStock
		}

		if err := s.bankRepo.UpdateStockQuantity(txCtx, stockName, -1); err != nil {
			return err
		}
		if err := s.walletRepo.UpdateStockQuantity(txCtx, walletID, stockName, 1); err != nil {
			return err
		}

		return s.logOperation(txCtx, domain.OperationTypeBuy, walletID, stockName)
	})
}

// SellStock executes a sale of a single stock from the wallet back to the bank.
func (s *ExchangeService) SellStock(ctx context.Context, walletID, stockName string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		_, err := s.ensureTradingPrerequisites(txCtx, walletID, stockName)
		if err != nil {
			return err
		}

		walletQty, err := s.walletRepo.GetStockQuantity(txCtx, walletID, stockName)
		if err != nil && !errors.Is(err, domain.ErrStockNotFound) {
			return err
		}
		if walletQty <= 0 {
			return domain.ErrInsufficientStock
		}

		if err := s.walletRepo.UpdateStockQuantity(txCtx, walletID, stockName, -1); err != nil {
			return err
		}
		if err := s.bankRepo.UpdateStockQuantity(txCtx, stockName, 1); err != nil {
			return err
		}

		return s.logOperation(txCtx, domain.OperationTypeSell, walletID, stockName)
	})
}

// ensureTradingPrerequisites validates if the bank has the stock and if the wallet exists, creating the wallet if it doesn't.
func (s *ExchangeService) ensureTradingPrerequisites(ctx context.Context, walletID, stockName string) (int, error) {
	bankQty, err := s.bankRepo.GetStockQuantity(ctx, stockName)
	if err != nil {
		return 0, err
	}

	_, err = s.walletRepo.GetWallet(ctx, walletID)
	if err != nil {
		if errors.Is(err, domain.ErrWalletNotFound) {
			if createErr := s.walletRepo.CreateWallet(ctx, walletID); createErr != nil {
				return 0, createErr
			}
		} else {
			return 0, err
		}
	}

	return bankQty, nil
}

// logOperation adds an entry to the audit log for a given operation.
func (s *ExchangeService) logOperation(ctx context.Context, op domain.OperationType, walletID, stockName string) error {
	return s.auditRepo.Add(ctx, domain.LogEntry{
		Type:      op,
		WalletID:  walletID,
		StockName: stockName,
	})
}
