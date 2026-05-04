package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/glitaa/stock-exchange/internal/domain"
)

// ErrorResponse represents a standard API error format.
type ErrorResponse struct {
	Error string `json:"error"`
}

// respondWithError maps domain errors to the correct HTTP status codes.
func respondWithError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError

	switch {
	case errors.Is(err, domain.ErrStockNotFound), errors.Is(err, domain.ErrWalletNotFound):
		statusCode = http.StatusNotFound
	case errors.Is(err, domain.ErrInsufficientStock), errors.Is(err, domain.ErrInvalidOperation):
		statusCode = http.StatusBadRequest
	}

	respondWithJSON(w, statusCode, ErrorResponse{Error: err.Error()})
}

// respondWithJSON is a helper to write JSON responses to the client.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to marshal response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
