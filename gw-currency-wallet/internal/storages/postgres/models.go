package postgres

import (
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"gorm.io/gorm"
)

func NewModels(db *gorm.DB) *storages.Models {
	return &storages.Models{
		Users:   NewUserModel(db),
		Wallets: NewWalletModel(db),
	}
}
