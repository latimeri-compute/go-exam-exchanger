package grpcclient

import (
	"context"

	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ExchangeResponse struct {
	Currency string  `json:"currency"`
	Rate     float32 `json:"rate"`
}

func NewClient(address string) (pb.ExchangeServiceClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}
	client := pb.NewExchangeServiceClient(conn)

	return client, nil
}

func GetOnlyRates(c pb.ExchangeServiceClient, ctx context.Context) ([]ExchangeResponse, error) {
	var response []ExchangeResponse
	rates, err := c.GetExchangeRates(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}
	for key, rate := range rates.Rates {
		response = append(response, ExchangeResponse{
			Currency: key,
			Rate:     rate,
		})
	}

	return response, nil
}

func GetOnlySpecificRate(c pb.ExchangeServiceClient, ctx context.Context, from, to string) (float32, error) {
	in := &pb.CurrencyRequest{
		FromCurrency: from,
		ToCurrency:   to,
	}
	resp, err := c.GetExchangeRateForCurrency(ctx, in, grpc.EmptyCallOption{})
	if err != nil {
		return 0, err
	}

	return resp.Rate, nil
}
