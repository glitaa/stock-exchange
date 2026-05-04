package handler

import (
	"encoding/json"
	"net/http"

	"github.com/glitaa/stock-exchange/internal/domain"
	"github.com/glitaa/stock-exchange/internal/service"
)

// BankHandler handles HTTP requests related to the bank's inventory.
type BankHandler struct {
	bankService *service.BankService
}

// NewBankHandler creates a new instance of BankHandler.
func NewBankHandler(s *service.BankService) *BankHandler {
	return &BankHandler{bankService: s}
}

// GetStocks handles the GET request to retrieve all stocks in the bank.
func (h *BankHandler) GetStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := h.bankService.GetStocks(r.Context())
	if err != nil {
		respondWithError(w, err)
		return
	}

	if stocks == nil {
		stocks = []domain.Stock{}
	}

	respondWithJSON(w, http.StatusOK, stocks)
}

// SetStocks handles the POST request to overwrite the bank's inventory.
func (h *BankHandler) SetStocks(w http.ResponseWriter, r *http.Request) {
	var stocks []domain.Stock
	if err := json.NewDecoder(r.Body).Decode(&stocks); err != nil {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request payload"})
		return
	}

	err := h.bankService.SetStocks(r.Context(), stocks)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
