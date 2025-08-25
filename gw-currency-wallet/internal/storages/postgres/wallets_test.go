package postgres

import (
	"testing"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	if testing.Short() {
		t.Skip("пропуск интеграционных тестов")
	}

	tests := []struct {
		name string
		ID   uint
		want error
	}{
		{
			name: "существующий",
			ID:   1,
			want: nil,
		},
		{
			name: "несуществующий",
			ID:   8888,
			want: storages.ErrRecordNotFound,
		},
		{
			name: "нулевой id",
			ID:   0,
			want: storages.ErrRecordNotFound,
		},
	}

	db := newTestDB(t)
	model := NewWalletModel(db)
	setupDB(t, db)
	defer teardownDB(t, db)

	createOneWallet(t, db)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := model.GetBalance(test.ID)
			assert.ErrorIs(t, err, test.want)
		})
	}
}

func TestChangeBalance(t *testing.T) {
	if testing.Short() {
		t.Skip("пропуск интеграционных тестов")
	}

	tests := []struct {
		name       string
		id         uint
		amount     int
		currency   string
		wantWallet storages.Wallet
		wantErr    error
	}{
		{
			name:     "существующий",
			amount:   1000000,
			currency: "rub",
			wantWallet: storages.Wallet{
				RubBalance: 1000000,
			},
			wantErr: nil,
		},
		{
			name:     "существующий",
			amount:   1000000,
			currency: "eur",
			wantWallet: storages.Wallet{
				EurBalance: 1000000,
			},
			wantErr: nil,
		},
		{
			name:     "существующий",
			amount:   1000000,
			currency: "usd",
			wantWallet: storages.Wallet{
				UsdBalance: 1000000,
			},
			wantErr: nil,
		},
		{
			name:     "несуществующий",
			id:       8989898,
			amount:   900,
			currency: "rub",
			wantErr:  storages.ErrRecordNotFound,
		},
		{
			name:     "слишком большая сумма для снятия",
			currency: "rub",
			amount:   -99999999,
			wantErr:  storages.ErrLessThanZero,
		},
	}

	db := newTestDB(t)
	model := NewWalletModel(db)
	setupDB(t, db)
	defer teardownDB(t, db)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := createOneWallet(t, db)
			if test.id == 0 {
				test.id = wallet.ID
			}

			wallet, err := model.ChangeBalance(test.id, test.amount, test.currency)
			t.Log(wallet, err)
			test.wantWallet.ID = wallet.ID
			test.wantWallet.Model = wallet.Model
			assert.ErrorIs(t, err, test.wantErr)
			assert.Equal(t, test.wantWallet, wallet)
		})
	}
}

func TestExchangeBetweenCurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("пропуск интеграционных тестов")
	}
	tests := []struct {
		name         string
		id           uint
		amount       int
		fromCurrency string
		toCurrency   string
		rate         int
		startWallet  storages.Wallet
		wantWallet   storages.Wallet
		wantErr      error
	}{
		{
			name:         "норм",
			amount:       100,
			fromCurrency: "usd",
			toCurrency:   "rub",
			rate:         8002,
			startWallet: storages.Wallet{
				RubBalance: 0,
				UsdBalance: 100,
			},
			wantWallet: storages.Wallet{
				RubBalance: 8002,
				UsdBalance: 0,
			},
			wantErr: nil,
		},
		{
			name:         "слишком большая сумма для снятия",
			amount:       100,
			fromCurrency: "usd",
			toCurrency:   "rub",
			rate:         -8002,
			startWallet: storages.Wallet{
				RubBalance: 0,
				UsdBalance: 100,
			},
			wantWallet: storages.Wallet{},
			wantErr:    storages.ErrLessThanZero,
		},
		{
			name:         "несуществующий кошелёк",
			id:           90000,
			amount:       900,
			fromCurrency: "rub",
			toCurrency:   "usd",
			// wantWallet:   storages.Wallet{},
			wantErr: storages.ErrRecordNotFound,
		},
	}

	db := newTestDB(t)
	model := NewWalletModel(db)
	setupDB(t, db)
	defer teardownDB(t, db)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := createOneWallet(t, db)
			if test.id == 0 {
				test.id = wallet.ID
			}
			wallet, err := model.ChangeBalance(test.id, int(test.startWallet.UsdBalance), "usd")
			if err != nil {
				t.Fatal(err)
			}

			wallet, err = model.ExchangeBetweenCurrency(test.id, test.amount, test.rate, test.fromCurrency, test.toCurrency)
			test.wantWallet.ID = wallet.ID
			test.wantWallet.Model = wallet.Model
			assert.ErrorIs(t, err, test.wantErr)
			assert.Equal(t, test.wantWallet, wallet)
		})
	}
}
