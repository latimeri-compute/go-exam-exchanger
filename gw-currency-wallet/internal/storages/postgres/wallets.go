package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
	column := fmt.Sprintf("%s_balance", strings.ToLower(currency))

	err := m.DB.Transaction(func(tx *gorm.DB) error {
		var balance int
		if err := tx.Model(wallet).Where("id = ?", wallet.ID).Select(column).First(&balance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return storages.ErrRecordNotFound
			}
			return err
		}

		newBalance := amount + balance
		if newBalance < 0 {
			return storages.ErrLessThanZero
		}

		if err := tx.Model(wallet).
			Where("id = ?", id).
			Update(column, newBalance).
			Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})

	if err != nil {
		return storages.Wallet{}, err
	}
	return *wallet, nil
}

// func (m *WalletModel) ExchangeBetweenCurrency(id uint, amount int, fromCurrency string, toCurrency string) (storages.Wallet, error) {
// 	wallet := &storages.Wallet{
// 		Model: gorm.Model{
// 			ID: id,
// 		},
// 	}
// 	fromColumn := fmt.Sprintf("%s_balance", strings.ToLower(fromCurrency))
// 	toColumn := fmt.Sprintf("%s_balance", strings.ToLower(toCurrency))

// 	err := m.DB.Transaction(func(tx *gorm.DB) error {
// 		var newFromBalance int
// 		if err := tx.Model(wallet).Where("id = ?", wallet.ID).Select(fromColumn).First(&newFromBalance).; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return storages.ErrRecordNotFound
// 			}
// 			return err
// 		}

// 		newFromBalance = newFromBalance - amount
// 		if newFromBalance < 0 {
// 			return storages.ErrLessThanZero
// 		}
// 		newToBalance :=

// 		if err := tx.Model(wallet).
// 			Where("id = ?", id).
// 			Update(toColumn, newBalance).
// 			Error; err != nil {
// 			return err
// 		}
// 		return nil
// 	}, &sql.TxOptions{
// 		Isolation: sql.LevelReadCommitted,
// 	})

// 	if err != nil {
// 		return storages.Wallet{}, err
// 	}
// 	return *wallet, nil
// }
