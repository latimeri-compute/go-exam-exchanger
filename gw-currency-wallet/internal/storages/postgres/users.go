package postgres

import (
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"gorm.io/gorm"
)

type UserModel struct {
	DB *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		DB: db,
	}
}

func (m *UserModel) CreateUser(user *storages.User) error {
	newWallet := &storages.Wallet{
		UsdBalance: 0,
		RubBalance: 0,
		EurBalance: 0,
	}

	err := m.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newWallet).Error; err != nil {
			return err
		}

		user.Wallet.ID = newWallet.ID
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (m *UserModel) FindUser(user *storages.User) error {
	err := m.DB.First(user).Error
	return err
}
