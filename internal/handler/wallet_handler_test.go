package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glitaa/stock-exchange/internal/domain"
	"github.com/glitaa/stock-exchange/internal/domain/mocks"
	"github.com/glitaa/stock-exchange/internal/service"
	"go.uber.org/mock/gomock"
)

func TestWalletHandler_GetWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWalletRepository(ctrl)
	svc := service.NewWalletService(mockRepo)
	h := NewWalletHandler(svc)

	mockRepo.EXPECT().GetWallet(gomock.Any(), "w1").Return(domain.Wallet{ID: "w1", Stocks: []domain.Stock{}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/wallets/w1", nil)
	req.SetPathValue("wallet_id", "w1")
	rr := httptest.NewRecorder()
	h.GetWallet(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var w domain.Wallet
	json.NewDecoder(rr.Body).Decode(&w)
	if w.ID != "w1" {
		t.Errorf("expected wallet w1, got %v", w.ID)
	}
}

func TestWalletHandler_GetWalletStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWalletRepository(ctrl)
	svc := service.NewWalletService(mockRepo)
	h := NewWalletHandler(svc)

	mockRepo.EXPECT().GetWallet(gomock.Any(), "w1").Return(domain.Wallet{ID: "w1"}, nil)
	mockRepo.EXPECT().GetStockQuantity(gomock.Any(), "w1", "AAPL").Return(5, nil)

	req := httptest.NewRequest(http.MethodGet, "/wallets/w1/stocks/AAPL", nil)
	req.SetPathValue("wallet_id", "w1")
	req.SetPathValue("stock_name", "AAPL")
	rr := httptest.NewRecorder()
	h.GetWalletStock(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var qty int
	json.NewDecoder(rr.Body).Decode(&qty)
	if qty != 5 {
		t.Errorf("expected 5, got %d", qty)
	}
}
