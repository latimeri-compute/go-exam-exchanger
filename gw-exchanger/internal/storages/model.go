package storages

import (
	"time"
)

type Valute struct {
	ID       int    `gorm:"primaryKey"`
	Code     string `gorm:"type:VARCHAR(6) NOT NULL;"`
	FullName string `gorm:"type:VARCHAR(255) NOT NULL;"`
}

type Exchange struct {
	FromValuteID int       `gorm:"type:int not null;"`
	ToValuteID   int       `gorm:"type:int not null;"`
	Rate         uint64    `gorm:"type:BIGINT NOT NULL;"`
	RateID       int       `gorm:"primaryKey"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	FromValute Valute `gorm:"foreignKey:FromValuteID;references:ID"`
	ToValute   Valute `gorm:"foreignKey:ToValuteID;references:ID"`
}

type ReturnExchanges struct {
	FromValuteCode string
	ToValuteCode   string
	Rate           uint64
}
