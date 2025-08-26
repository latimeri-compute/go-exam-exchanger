package mock_storages

import (
	"strings"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"gorm.io/gorm"
)

var ValidWallet = storages.Wallet{
	Model: gorm.Model{
		ID:        1,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	},
	UsdBalance: 0,
	RubBalance: 0,
	EurBalance: 0,
}

type MockWallets struct{}

func NewMockWallets() *MockWallets {
	return &MockWallets{}
}

func (n *MockWallets) GetBalance(id uint) (storages.Wallet, error) {
	if id == 1 {
		return ValidWallet, nil
	} else {
		return storages.Wallet{}, storages.ErrRecordNotFound
	}
}

func (n *MockWallets) ChangeBalance(id uint, amount int, currency string) (storages.Wallet, error) {
	var wallet storages.Wallet
	if id != 1 {
		return wallet, storages.ErrRecordNotFound
	}
	if amount < -10000 {
		return wallet, storages.ErrLessThanZero
	}
	wallet.ID = 1
	wallet.RubBalance = 10000
	wallet.UsdBalance = 10000
	wallet.EurBalance = 10000

	switch strings.ToLower(currency) {
	case "rub":
		wallet.RubBalance += int64(amount)
	case "usd":
		wallet.UsdBalance += int64(amount)
	case "eur":
		wallet.EurBalance += int64(amount)
	default:
		panic("ChangeBalance: валюта не поддерживается")
	}

	return wallet, nil
}

func (n *MockWallets) ExchangeBetweenCurrency(id uint, amount int, rate int, fromCurrency string, toCurrency string) (storages.Wallet, error) {
	var wallet storages.Wallet
	if id != 1 {
		return wallet, storages.ErrRecordNotFound
	}
	if amount > 10000 {
		return wallet, storages.ErrLessThanZero
	}

	wallet.ID = 1
	wallet.RubBalance = 10000
	wallet.UsdBalance = 10000
	wallet.EurBalance = 10000

	// uhhhh в принципе я не думаю, что в моках нужно учитывать курс;
	// он только усложняет проверку родительских методов
	switch strings.ToLower(toCurrency) {
	case "rub":
		wallet.RubBalance += int64(amount)
	case "usd":
		wallet.UsdBalance += int64(amount)
	case "eur":
		wallet.EurBalance += int64(amount)
	default:
		panic("ExchangeBetweenCurrency: валюта не поддерживается")
	}
	switch strings.ToLower(fromCurrency) {
	case "rub":
		wallet.RubBalance -= int64(amount)
	case "usd":
		wallet.UsdBalance -= int64(amount)
	case "eur":
		wallet.EurBalance -= int64(amount)
	default:
		panic("ExchangeBetweenCurrency: валюта не поддерживается")
	}

	return wallet, nil
}
