package storages

import (
	"time"
)

type Valute struct {
	ID       int    `gorm:"primaryKey"`
	Code     string `gorm:"type:VARCHAR(6) NOT NULL;"`
	FullName string `gorm:"type:VARCHAR(255) NOT NULL;"`
}

type ExchangeRate struct {
	FromValuteID int       `gorm:"type:int not null;foreignKey:ValuteID;"`
	ToValuteID   int       `gorm:"type:int not null;foreignKey:ValuteID;"`
	Rate         uint64    `gorm:"type:BIGINT NOT NULL;"`
	RateId       int       `gorm:"primaryKey"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

type ReturnExchanges struct {
	FromValuteCode string
	ToValuteCode   string
	Rate           uint64
}

type Models struct {
	ExchangerModel ExchangerModelInterface
}
