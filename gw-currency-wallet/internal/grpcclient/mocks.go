package grpcclient

import (
	"context"
	"fmt"

	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
)

var Rates = map[string]float32{
	"rub->usd": 0.4,
	"usd->rub": 85,
	"rub->eur": 0.4,
	"eur->rub": 93,
	"usd->eur": 0.8,
	"eur->usd": 1.4,
}

type Handler struct {
	pb.UnimplementedExchangeServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetExchangeRates(ctx context.Context, in *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	var res pb.ExchangeRatesResponse
	res.Rates = Rates
	return &res, nil
}
func (h *Handler) GetExchangeRateForCurrency(ctx context.Context, in *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	var res pb.ExchangeRateResponse

	key := fmt.Sprintf("%s->%s", in.FromCurrency, in.ToCurrency)
	res.Rate = Rates[key]
	return &res, nil
}
