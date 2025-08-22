package postgres

import (
	"errors"

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
	wal := &storages.Wallet{}
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(wal).Error; err != nil {
			return err
		}

		user.WalletID = wal.ID
		if err := tx.Create(user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return storages.ErrRecordExists
			}
			return err
		}
		return nil
	})

	return err
}

func (m *UserModel) FindUser(user *storages.User) error {
	err := m.DB.First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return storages.ErrRecordNotFound
	}
	return err
}
