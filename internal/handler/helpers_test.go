package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glitaa/stock-exchange/internal/domain"
)

func TestRespondWithJSON(t *testing.T) {
	t.Run("successful JSON response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		payload := map[string]string{"message": "success"}
		
		respondWithJSON(recorder, http.StatusOK, payload)

		if recorder.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
		}

		if recorder.Header().Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", recorder.Header().Get("Content-Type"))
		}

		var response map[string]string
		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		if response["message"] != "success" {
			t.Errorf("expected message 'success', got '%s'", response["message"])
		}
	})

	t.Run("JSON marshaling failure", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		// Channels cannot be marshaled to JSON
		payload := make(chan int)
		
		respondWithJSON(recorder, http.StatusOK, payload)

		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
		}

		expectedBody := "{\"error\":\"failed to marshal response\"}"
		if recorder.Body.String() != expectedBody {
			t.Errorf("expected body '%s', got '%s'", expectedBody, recorder.Body.String())
		}
	})
}

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "Not Found - Stock",
			err:            domain.ErrStockNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Not Found - Wallet",
			err:            domain.ErrWalletNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Bad Request - Insufficient Stock",
			err:            domain.ErrInsufficientStock,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Bad Request - Invalid Operation",
			err:            domain.ErrInvalidOperation,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Internal Server Error - Generic",
			err:            errors.New("some unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			respondWithError(recorder, tt.err)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}

			var response ErrorResponse
			if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if response.Error != tt.err.Error() {
				t.Errorf("expected error message '%s', got '%s'", tt.err.Error(), response.Error)
			}
		})
	}
}
