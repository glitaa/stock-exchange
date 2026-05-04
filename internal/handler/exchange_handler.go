package handler

import (
	"encoding/json"
	"net/http"

	"github.com/glitaa/stock-exchange/internal/service"
)

// ExchangeHandler handles HTTP requests related to trading.
type ExchangeHandler struct {
	exchangeService *service.ExchangeService
}

// NewExchangeHandler creates a new instance of ExchangeHandler.
func NewExchangeHandler(s *service.ExchangeService) *ExchangeHandler {
	return &ExchangeHandler{exchangeService: s}
}

// tradeRequest defines the expected JSON payload for buy and sell operations.
type tradeRequest struct {
	WalletID  string `json:"wallet_id"`
	StockName string `json:"stock_name"`
}

// BuyStock handles the POST request to buy a stock.
func (h *ExchangeHandler) BuyStock(w http.ResponseWriter, r *http.Request) {
	var req tradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request payload"})
		return
	}

	err := h.exchangeService.BuyStock(r.Context(), req.WalletID, req.StockName)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// SellStock handles the POST request to sell a stock.
func (h *ExchangeHandler) SellStock(w http.ResponseWriter, r *http.Request) {
	var req tradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request payload"})
		return
	}

	err := h.exchangeService.SellStock(r.Context(), req.WalletID, req.StockName)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
