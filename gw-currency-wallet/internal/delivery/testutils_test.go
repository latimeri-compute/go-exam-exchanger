package delivery

import (
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/brocker"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/cache"
	mock_storages "github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages/mocks"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"go.uber.org/zap"
)

func NewTestHandler(jwtString string, grpcClient pb.ExchangeServiceClient, messenger *brocker.Producer) *Handler {
	c := cache.New[string, []ExchangeCachedItem]()
	return &Handler{
		Models:         mock_storages.NewMockModels(),
		Logger:         zap.NewNop().Sugar(),
		messenger:      messenger,
		exchangeCache:  c,
		ExchangeClient: grpcClient,
		JWTsource:      jwtString,
	}
}
