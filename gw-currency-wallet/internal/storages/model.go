package storages

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("Record not found")
	ErrRecordExists   = errors.New("Record already exists")
)

type User struct {
	gorm.Model

	Email        string `gorm:"type:varchar(255);not null;unique;"`
	PasswordHash []byte `gorm:"bytea;not null;"`
	JWTToken     string `gotm:"varchar(255);not null;"`

	WalletID uint `gorm:"column:wallet_id;foreignKey:wallet_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Wallet   Wallet
}

type Wallet struct {
	gorm.Model

	UsdBalance int64 `gorm:"type:bigint;not null;default:0;"`
	EurBalance int64 `gorm:"type:bigint;not null;default:0;"`
	RubBalance int64 `gorm:"type:bigint;not null;default:0;"`
}

type Models struct {
	Users   UserModelInterface
	Wallets WalletModelInterface
}
