package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	wallet := storages.Wallet{
		Model: gorm.Model{
			ID: id,
		},
	}
	err := m.DB.Select("usd_balance", "rub_balance", "eur_balance").Where("id = ?", wallet.ID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return storages.Wallet{}, storages.ErrRecordNotFound
		}
		return storages.Wallet{}, err
	}

	return wallet, nil
}

// последние две цифры это десятичные
func (m *WalletModel) ChangeBalance(id uint, amount int, currency string) (storages.Wallet, error) {
	wallet := &storages.Wallet{
		Model: gorm.Model{
			ID: id,
		},
	}
	column := fmt.Sprintf("%s_balance", strings.ToLower(currency))

	err := m.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&wallet).
			Clauses(clause.Returning{}).
			Where("id = ?", wallet.ID).
			Update(column, gorm.Expr(fmt.Sprintf("%s + ?", column), amount))

		err := res.Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return storages.ErrRecordNotFound
			}
			return err
		}

		if res.RowsAffected == 0 {
			return storages.ErrRecordNotFound
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

func (m *WalletModel) ExchangeBetweenCurrency(id uint, amount int, rate int, fromCurrency string, toCurrency string) (storages.Wallet, error) {
	wallet := storages.Wallet{
		Model: gorm.Model{
			ID: id,
		},
	}
	fromColumn := fmt.Sprintf("%s_balance", strings.ToLower(fromCurrency))
	toColumn := fmt.Sprintf("%s_balance", strings.ToLower(toCurrency))

	amount = utils.Abs(amount)

	err := m.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&wallet).
			Clauses(clause.Returning{}).
			Where("id = ?", wallet.ID).
			Updates(map[string]any{
				fromColumn: gorm.Expr(fmt.Sprintf("%s - ?", fromColumn), amount),
				toColumn:   gorm.Expr(fmt.Sprintf("%s + ?", toColumn), amount*rate/100),
			})

		err := res.Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return storages.ErrRecordNotFound
			}
			return err
		}

		if res.RowsAffected == 0 {
			return storages.ErrRecordNotFound
		}

		return nil

	}, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})

	if err != nil {
		return storages.Wallet{}, err
	}
	return wallet, nil
}
