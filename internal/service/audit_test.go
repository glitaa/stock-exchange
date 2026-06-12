package service

import (
	"context"
	"testing"

	"github.com/glitaa/stock-exchange/internal/domain"
	"github.com/glitaa/stock-exchange/internal/domain/mocks"
	"go.uber.org/mock/gomock"
)

func TestAuditService_GetLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuditLogRepository(ctrl)
	svc := NewAuditService(mockRepo)
	ctx := context.Background()

	expectedLogs := []domain.LogEntry{{Type: domain.OperationTypeBuy, WalletID: "w1", StockName: "AAPL"}}
	mockRepo.EXPECT().GetAll(ctx).Return(expectedLogs, nil)

	logs, err := svc.GetLog(ctx)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}
	if len(logs) != 1 || logs[0].Type != domain.OperationTypeBuy {
		t.Errorf("expected Buy log, got %v", logs)
	}
}
