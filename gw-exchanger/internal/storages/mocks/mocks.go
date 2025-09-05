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

var ExchangeReturns = []storages.ReturnExchanges{
	{
		FromValuteCode: "rub",
		ToValuteCode:   "usd",
		Rate:           560000,
	},
	{
		FromValuteCode: "rub",
		ToValuteCode:   "eur",
		Rate:           1000000,
	},
	{
		FromValuteCode: "usd",
		ToValuteCode:   "eur",
		Rate:           9000,
	},
	{
		FromValuteCode: "usd",
		ToValuteCode:   "rub",
		Rate:           37,
	},
}

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
	{
		FromValuteID: 2,
		ToValuteID:   1,
		Rate:         37,
		RateID:       4,
		UpdatedAt:    time.Now(),
		FromValute:   ValuteUSD,
		ToValute:     ValuteRub,
	},
}

type MockExchange struct {
}

func NewExchange() *MockExchange {
	return &MockExchange{}
}

func (m *MockExchange) GetAll() ([]storages.ReturnExchanges, error) {
	return ExchangeReturns, nil
}

func (m *MockExchange) GetRateBetween(fromValute, toValute string) (storages.ReturnExchanges, error) {
	var res storages.ReturnExchanges
	for _, e := range ExchangeReturns {
		if fromValute == e.FromValuteCode && toValute == e.ToValuteCode {
			res = storages.ReturnExchanges{
				FromValuteCode: e.FromValuteCode,
				ToValuteCode:   e.ToValuteCode,
				Rate:           e.Rate,
			}
		}
	}
	if res.FromValuteCode == "" {
		return storages.ReturnExchanges{}, storages.ErrNotFound
	}
	return res, nil
}
