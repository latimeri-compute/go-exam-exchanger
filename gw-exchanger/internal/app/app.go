package app

import (
	delivery "gitgub.com/latimeri-compute/gw-exchanger/internal/delivery/response"
	"gitgub.com/latimeri-compute/gw-exchanger/internal/storages"
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
