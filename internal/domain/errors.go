package domain

import "errors"

var (
	ErrStockNotFound     = errors.New("stock not found")
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidOperation  = errors.New("invalid operation")
)
