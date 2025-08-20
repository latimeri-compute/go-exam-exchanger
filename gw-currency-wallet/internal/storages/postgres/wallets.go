package postgres

import (
	"fmt"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"gorm.io/gorm"
)

type WalletModel struct {
	DB *gorm.DB
}

func NewWalletModel(db *gorm.DB) *WalletModel {
	return &WalletModel{
		DB: db,
	}
}

func (m *WalletModel) GetBalance(id uint) (storages.Wallet, error) {
	wallet := &storages.Wallet{
		Model: gorm.Model{
			ID: id,
		},
	}
	err := m.DB.Select("usd_balance", "rub_balance", "eur_balance").First(wallet).Error

	return *wallet, err
}

// последние четыре цифры это десятичные
func (m *WalletModel) ChangeBalance(id uint, amount int, currency string) (storages.Wallet, error) {
	wallet := &storages.Wallet{
		Model: gorm.Model{
			ID: id,
		},
	}
	column := fmt.Sprintf("%s_balance", currency)

	err := m.DB.Transaction(func(tx *gorm.DB) error {
		var balance *int
		if err := tx.Model(wallet).Select(column).First(balance).Error; err != nil {
			return err
		}

		newBalance := amount + *balance
		if err := tx.Model(&storages.Wallet{}).
			Where("id = ? AND (balance + ?) >= 0", id, amount).
			Update(column, newBalance).
			Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return storages.Wallet{}, err
	}
	return *wallet, nil
}
