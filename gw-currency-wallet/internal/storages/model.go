package storages

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email        string `gorm:"type:VARCHAR(255) NOT NULL;unique;"`
	PasswordHash []byte `gorm:"bytea NOT NULL;"`

	Wallet Wallet `gorm:"foreignKey:WalletID;references:ID"`
}

type Wallet struct {
	gorm.Model

	UsdBalance uint64 `gorm:"type:BIGINT NOT NULL;"`
	EurBalance uint64 `gorm:"type:BIGINT NOT NULL;"`
	RubBalance uint64 `gorm:"type:BIGINT NOT NULL;"`
}

type Models struct {
	Users   UserModelInterface
	Wallets WalletModelInterface
}
