package delivery

import (
	"context"

	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/"
	"go.uber.org/zap"
)

type Handler struct {
	pb.ExchangeServiceClient
	pb.UnimplementedExchangeServiceServer

	logger *zap.Logger
}

func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// Получение курсов обмена всех валют
func (h *Handler) GetExchangeRates(ctx context.Context, in *pb.Empty) (*pb.ExchangeRatesResponse, error) {

	return nil, nil
}

// Получение курса обмена для конкретной валюты
func (h *Handler) GetExchangeRateForCurrency(ctx context.Context, in *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	return nil, nil
}
