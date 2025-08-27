package delivery

import (
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/cache"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"go.uber.org/zap"
)

type Handler struct {
	Models         *storages.Models
	Logger         *zap.SugaredLogger
	ExchangeClient pb.ExchangeServiceClient
	exchangeCache  *cache.Cache[string, []ExchangeCachedItem]
	JWTsource      string
}

func NewHandler(m *storages.Models, logger *zap.SugaredLogger, exchangeClient pb.ExchangeServiceClient, jwt string) *Handler {
	return &Handler{
		Models:         m,
		Logger:         logger,
		JWTsource:      jwt,
		ExchangeClient: exchangeClient,
		exchangeCache:  cache.New[string, []ExchangeCachedItem](),
	}
}

type ExchangeCachedItem struct {
	Rate         float32 `json:"rate"`
	ToCurrency   string  `json:"to_currency"`
	FromCurrency string  `json:"from_currency"`
}
