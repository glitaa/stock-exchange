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

// bankResponse defines the expected JSON structure for bank endpoints.
type bankPayload struct {
	Stocks []domain.Stock `json:"stocks"`
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

	response := bankPayload{Stocks: stocks}
	respondWithJSON(w, http.StatusOK, response)
}

// SetStocks handles the POST request to overwrite the bank's inventory.
func (h *BankHandler) SetStocks(w http.ResponseWriter, r *http.Request) {
	var payload bankPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request payload"})
		return
	}

	err := h.bankService.SetStocks(r.Context(), payload.Stocks)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
