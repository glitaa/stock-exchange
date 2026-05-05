package handler

import (
	"net/http"

	"github.com/glitaa/stock-exchange/internal/domain"
	"github.com/glitaa/stock-exchange/internal/service"
)

// WalletHandler handles HTTP requests related to wallets.
type WalletHandler struct {
	walletService *service.WalletService
}

// NewWalletHandler creates a new instance of WalletHandler.
func NewWalletHandler(s *service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: s}
}

// GetWallet handles the GET request to retrieve a wallet's state.
func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("wallet_id")
	if walletID == "" {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "wallet_id is required"})
		return
	}

	wallet, err := h.walletService.GetWallet(r.Context(), walletID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	if wallet.Stocks == nil {
		wallet.Stocks = []domain.Stock{}
	}

	respondWithJSON(w, http.StatusOK, wallet)
}

// GetWalletStock handles the GET request to retrieve the quantity of a specific stock in a wallet.
func (h *WalletHandler) GetWalletStock(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("wallet_id")
	stockName := r.PathValue("stock_name")

	if walletID == "" || stockName == "" {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "wallet_id and stock_name are required"})
		return
	}

	quantity, err := h.walletService.GetStockQuantity(r.Context(), walletID, stockName)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, quantity)
}
