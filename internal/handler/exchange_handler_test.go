package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glitaa/stock-exchange/internal/domain"
	"github.com/glitaa/stock-exchange/internal/domain/mocks"
	"github.com/glitaa/stock-exchange/internal/service"
	"go.uber.org/mock/gomock"
)

func TestExchangeHandler_Trade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mocks.NewMockWalletRepository(ctrl)
	mockBankRepo := mocks.NewMockBankRepository(ctrl)
	mockAuditRepo := mocks.NewMockAuditLogRepository(ctrl)
	mockTxManager := mocks.NewMockTxManager(ctrl)

	svc := service.NewExchangeService(mockWalletRepo, mockBankRepo, mockAuditRepo, mockTxManager)
	h := NewExchangeHandler(svc)

	t.Run("successful buy", func(t *testing.T) {
		mockWalletRepo.EXPECT().GetWallet(gomock.Any(), "w1").Return(domain.Wallet{}, nil)
		mockTxManager.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
		mockBankRepo.EXPECT().GetStockQuantity(gomock.Any(), "AAPL").Return(10, nil)
		mockBankRepo.EXPECT().UpdateStockQuantity(gomock.Any(), "AAPL", -1).Return(nil)
		mockWalletRepo.EXPECT().UpdateStockQuantity(gomock.Any(), "w1", "AAPL", 1).Return(nil)
		mockAuditRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)

		body := []byte(`{"type":"buy"}`)
		req := httptest.NewRequest(http.MethodPost, "/wallets/w1/stocks/AAPL", bytes.NewReader(body))
		req.SetPathValue("wallet_id", "w1")
		req.SetPathValue("stock_name", "AAPL")
		rr := httptest.NewRecorder()

		h.Trade(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}
	})

	t.Run("bad request - invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/wallets/w1/stocks/AAPL", bytes.NewReader([]byte(`{"type":"invalid"}`)))
		req.SetPathValue("wallet_id", "w1")
		req.SetPathValue("stock_name", "AAPL")
		rr := httptest.NewRecorder()

		h.Trade(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rr.Code)
		}
	})
}
