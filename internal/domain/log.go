package domain

// OperationType defines the type of a stock market transaction.
type OperationType string

const (
	OperationTypeBuy  OperationType = "buy"
	OperationTypeSell OperationType = "sell"
)

// LogEntry represents a single record in the audit log.
type LogEntry struct {
	Type      OperationType `json:"type"`
	WalletID  string        `json:"wallet_id"`
	StockName string        `json:"stock_name"`
}
