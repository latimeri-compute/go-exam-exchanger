package app

import (
	delivery "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/internal/delivery/response"
	"github.com/latimeri-compute/go-exam-exchanger/proto-exchange/internal/storages"
	"go.uber.org/zap"
)

type App struct {
	Handler *delivery.Handler
	Models  *storages.ExchangerModelInterface
}

func New(logger *zap.Logger) *App {
	return &App{

		Handler: delivery.NewHandler(logger),
	}
}
