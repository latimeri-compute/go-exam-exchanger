package app

import (
	"github.com/go-chi/chi"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"go.uber.org/zap"
)

type App struct {
	Router *chi.Mux
	logger *zap.Logger
	Models *storages.Models
}
