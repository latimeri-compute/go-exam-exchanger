package delivery

import (
	mock_storages "github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages/mocks"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"go.uber.org/zap"
)

func NewTestHandler(jwtString string, grpcClient pb.ExchangeServiceClient) *Handler {
	return &Handler{
		Models:         mock_storages.NewMockModels(),
		Logger:         zap.NewNop().Sugar(),
		ExchangeClient: grpcClient,
		JWTsource:      jwtString,
	}
}
