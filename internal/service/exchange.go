package service

import (
	"context"

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

// ensureWalletExists is a helper that creates a wallet if it doesn't already exist.
func (s *ExchangeService) ensureWalletExists(ctx context.Context, walletID string) error {
	_, err := s.walletRepo.GetWallet(ctx, walletID)
	if err != nil {
		if err == domain.ErrWalletNotFound {
			return s.walletRepo.CreateWallet(ctx, walletID)
		}
		return err
	}
	return nil
}

// BuyStock executes a purchase of a single stock from the bank to the wallet.
func (s *ExchangeService) BuyStock(ctx context.Context, walletID, stockName string) error {
	if err := s.ensureWalletExists(ctx, walletID); err != nil {
		return err
	}

	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		bankQty, err := s.bankRepo.GetStockQuantity(txCtx, stockName)
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
	if err := s.ensureWalletExists(ctx, walletID); err != nil {
		return err
	}

	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		walletQty, err := s.walletRepo.GetStockQuantity(txCtx, walletID, stockName)
		if err != nil {
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

// logOperation adds an entry to the audit log for a given operation.
func (s *ExchangeService) logOperation(ctx context.Context, op domain.OperationType, walletID, stockName string) error {
	return s.auditRepo.Add(ctx, domain.LogEntry{
		Type:      op,
		WalletID:  walletID,
		StockName: stockName,
	})
}
