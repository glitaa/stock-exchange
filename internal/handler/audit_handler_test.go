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

func TestAuditHandler_GetLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuditLogRepository(ctrl)
	svc := service.NewAuditService(mockRepo)
	h := NewAuditHandler(svc)

	mockRepo.EXPECT().GetAll(gomock.Any()).Return([]domain.LogEntry{{Type: domain.OperationTypeBuy, WalletID: "w1", StockName: "AAPL"}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/log", nil)
	rr := httptest.NewRecorder()
	h.GetLog(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var response auditResponse
	json.NewDecoder(rr.Body).Decode(&response)
	if len(response.Log) != 1 {
		t.Errorf("expected 1 log, got %d", len(response.Log))
	}
}
