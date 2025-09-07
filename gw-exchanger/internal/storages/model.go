package storages

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("Record not found")

type Valute struct {
	ID       int    `gorm:"primaryKey"`
	Code     string `gorm:"type:VARCHAR(6) NOT NULL;"`
	FullName string `gorm:"type:VARCHAR(255) NOT NULL;"`
}

type Exchange struct {
	FromValuteID int       `gorm:"type:int not null;index:unique_combo;"`
	ToValuteID   int       `gorm:"type:int not null;index:unique_combo;"`
	Rate         uint64    `gorm:"type:BIGINT NOT NULL;"`
	RateID       int       `gorm:"primaryKey"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	FromValute Valute `gorm:"foreignKey:FromValuteID;references:ID"`
	ToValute   Valute `gorm:"foreignKey:ToValuteID;references:ID"`
}

type ReturnExchanges struct {
	FromValuteCode string `gorm:"from_valute_code"`
	ToValuteCode   string `gorm:"to_valute_code"`
	Rate           uint64
}
