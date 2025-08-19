package mocks

import (
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages"
)

var (
	ValuteRub = storages.Valute{
		ID:       1,
		Code:     "rub",
		FullName: "Russian Ruble",
	}
	ValuteUSD = storages.Valute{
		ID:       2,
		Code:     "usd",
		FullName: "United States Dollar",
	}
	ValuteEUR = storages.Valute{
		ID:       3,
		Code:     "eur",
		FullName: "Euro",
	}
)

var Exchanges = []storages.Exchange{
	{
		FromValuteID: 1,
		ToValuteID:   2,
		Rate:         560000,
		RateID:       1,
		UpdatedAt:    time.Now(),
		FromValute:   ValuteRub,
		ToValute:     ValuteUSD,
	},
	{
		FromValuteID: 1,
		ToValuteID:   3,
		Rate:         1000000,
		RateID:       2,
		UpdatedAt:    time.Now(),
		FromValute:   ValuteRub,
		ToValute:     ValuteEUR,
	},
	{
		FromValuteID: 2,
		ToValuteID:   3,
		Rate:         9000,
		RateID:       3,
		UpdatedAt:    time.Now(),
		FromValute:   ValuteUSD,
		ToValute:     ValuteEUR,
	},
}

type MockExchange struct {
}

func NewExchange() *MockExchange {
	return &MockExchange{}
}

func (m *MockExchange) GetAll() ([]storages.Exchange, error) {
	return Exchanges, nil
}

func (m *MockExchange) GetRateBetween() {

}
