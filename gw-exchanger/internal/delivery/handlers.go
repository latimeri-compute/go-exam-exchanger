package delivery

import (
	"context"
	"fmt"

	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"go.uber.org/zap"
)

type Handler struct {
	pb.ExchangeServiceClient
	pb.UnimplementedExchangeServiceServer

	logger *zap.Logger
	db     storages.ExchangerModelInterface
}

func NewHandler(logger *zap.Logger, db storages.ExchangerModelInterface) *Handler {
	return &Handler{
		logger: logger,
		db:     db,
	}
}

// Получение курсов обмена всех валют
func (h *Handler) GetExchangeRates(ctx context.Context, in *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	var res pb.ExchangeRatesResponse
	rates := make(map[string]float32)

	exchange, err := h.db.GetAll()
	if err != nil {
		return &res, err
	}

	for _, e := range exchange {
		dir := fmt.Sprintf("%s->%s", e.FromValute.Code, e.ToValute.Code)
		rates[dir] = float32(e.Rate) / 10000
	}

	res.Rates = rates
	return &res, nil
}

// Получение курса обмена для конкретной валюты
func (h *Handler) GetExchangeRateForCurrency(ctx context.Context, in *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	return nil, nil
}
