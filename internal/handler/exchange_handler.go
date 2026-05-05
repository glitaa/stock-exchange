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

// tradeRequest defines the expected JSON payload for trade operations.
type tradeRequest struct {
	Type string `json:"type"` // "buy" or "sell"
}

// Trade handles the POST request to buy or sell a stock.
func (h *ExchangeHandler) Trade(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("wallet_id")
	stockName := r.PathValue("stock_name")

	if walletID == "" || stockName == "" {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "wallet_id and stock_name are required"})
		return
	}

	var req tradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request payload"})
		return
	}

	var err error
	switch req.Type {
	case "buy":
		err = h.exchangeService.BuyStock(r.Context(), walletID, stockName)
	case "sell":
		err = h.exchangeService.SellStock(r.Context(), walletID, stockName)
	default:
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid operation type, must be 'buy' or 'sell'"})
		return
	}

	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
