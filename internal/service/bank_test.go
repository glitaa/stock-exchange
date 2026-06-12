package service

import (
	"context"
	"testing"

	"github.com/glitaa/stock-exchange/internal/domain"
	"github.com/glitaa/stock-exchange/internal/domain/mocks"
	"go.uber.org/mock/gomock"
)

func TestBankService_GetStocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBankRepository(ctrl)
	svc := NewBankService(mockRepo)
	ctx := context.Background()

	expectedStocks := []domain.Stock{{Name: "AAPL", Quantity: 10}}
	mockRepo.EXPECT().GetStocks(ctx).Return(expectedStocks, nil)

	stocks, err := svc.GetStocks(ctx)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}
	if len(stocks) != 1 || stocks[0].Name != "AAPL" {
		t.Errorf("expected AAPL, got %v", stocks)
	}
}

func TestBankService_SetStocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBankRepository(ctrl)
	svc := NewBankService(mockRepo)
	ctx := context.Background()

	input := []domain.Stock{{Name: "AAPL", Quantity: 10}}
	mockRepo.EXPECT().SetStocks(ctx, input).Return(nil)

	err := svc.SetStocks(ctx, input)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}
}
