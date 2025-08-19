package delivery

import (
	"context"
	"errors"
	"fmt"

	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"gorm.io/gorm"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		h.logger.Error(err.Error())
		return &res, status.Error(codes.Internal, err.Error())
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
	var res pb.ExchangeRateResponse
	exchange, err := h.db.GetRateBetween(in.FromCurrency, in.ToCurrency)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Error(err.Error())
		}
		return &res, status.Error(codes.NotFound, err.Error())
	}

	res.FromCurrency = exchange.FromValute.Code
	res.ToCurrency = exchange.ToValute.Code
	res.Rate = float32(exchange.Rate) / 10000
	return &res, nil
}
