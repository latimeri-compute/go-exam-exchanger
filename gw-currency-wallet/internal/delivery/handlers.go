package delivery

import (
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"go.uber.org/zap"
)

type Handler struct {
	Models         *storages.Models
	Logger         *zap.SugaredLogger
	ExchangeClient pb.ExchangeServiceClient
	JWTsource      string
}

func NewHandler(m *storages.Models, logger *zap.SugaredLogger, jwt string) *Handler {
	return &Handler{
		Models:    m,
		Logger:    logger,
		JWTsource: jwt,
	}
}
