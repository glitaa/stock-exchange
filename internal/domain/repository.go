package domain

import "context"

// BankRepository handles operations related to the bank's stock inventory.
type BankRepository interface {
	GetStocks(ctx context.Context) ([]Stock, error)
	SetStocks(ctx context.Context, stocks []Stock) error
	GetStockQuantity(ctx context.Context, name string) (int, error)
	UpdateStockQuantity(ctx context.Context, name string, delta int) error
}

// WalletRepository handles operations related to user wallets.
type WalletRepository interface {
	GetWallet(ctx context.Context, walletID string) (Wallet, error)
	CreateWallet(ctx context.Context, walletID string) error
	GetStockQuantity(ctx context.Context, walletID string, stockName string) (int, error)
	UpdateStockQuantity(ctx context.Context, walletID string, stockName string, delta int) error
}

// AuditLogRepository handles operations related to the audit log.
type AuditLogRepository interface {
	Add(ctx context.Context, entry LogEntry) error
	GetAll(ctx context.Context) ([]LogEntry, error)
}

// TxManager defines a contract for executing operations within a database transaction.
type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
