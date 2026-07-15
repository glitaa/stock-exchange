package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glitaa/stock-exchange/internal/domain"
	"github.com/glitaa/stock-exchange/internal/domain/mocks"
	"github.com/glitaa/stock-exchange/internal/service"
	"go.uber.org/mock/gomock"
)

func TestBankHandler_GetStocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBankRepository(ctrl)
	bankSvc := service.NewBankService(mockRepo)
	h := NewBankHandler(bankSvc)

	mockRepo.EXPECT().GetStocks(gomock.Any()).Return([]domain.Stock{{Name: "AAPL", Quantity: 10}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/stocks", nil)
	rr := httptest.NewRecorder()
	h.GetStocks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var payload bankPayload
	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(payload.Stocks) != 1 || payload.Stocks[0].Name != "AAPL" {
		t.Errorf("expected AAPL stock, got %v", payload.Stocks)
	}
}

func TestBankHandler_SetStocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBankRepository(ctrl)
	bankSvc := service.NewBankService(mockRepo)
	h := NewBankHandler(bankSvc)

	mockRepo.EXPECT().SetStocks(gomock.Any(), gomock.Any()).Return(nil)

	body := []byte(`{"stocks":[{"name":"AAPL","quantity":10}]}`)
	req := httptest.NewRequest(http.MethodPost, "/stocks", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	h.SetStocks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
