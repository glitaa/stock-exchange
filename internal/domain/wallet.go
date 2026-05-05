package domain

// Wallet represents a user's wallet and the stocks they own.
type Wallet struct {
	ID     string  `json:"id"`
	Stocks []Stock `json:"stocks"`
}
