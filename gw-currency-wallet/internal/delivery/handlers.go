package delivery

import (
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/brocker"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/cache"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"go.uber.org/zap"
)

type Handler struct {
	Models         *storages.Models
	Logger         *zap.SugaredLogger
	ExchangeClient pb.ExchangeServiceClient
	messenger      *brocker.Producer
	exchangeCache  *cache.Cache[string, []ExchangeCachedItem]
	JWTsource      string
}
type ExchangeCachedItem struct {
	Rate         float32 `json:"rate"`
	ToCurrency   string  `json:"to_currency"`
	FromCurrency string  `json:"from_currency"`
}

func NewHandler(m *storages.Models, logger *zap.SugaredLogger, exchangeClient pb.ExchangeServiceClient,
	messenger *brocker.Producer, jwt string) *Handler {
	return &Handler{
		Models:         m,
		Logger:         logger,
		JWTsource:      jwt,
		ExchangeClient: exchangeClient,
		messenger:      messenger,
		exchangeCache:  cache.New[string, []ExchangeCachedItem](),
	}
}

func (h *Handler) UnpackJSON() {

}
