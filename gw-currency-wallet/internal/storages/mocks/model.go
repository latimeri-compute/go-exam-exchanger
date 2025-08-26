package mock_storages

import "github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"

func NewMockModels() *storages.Models {
	return &storages.Models{
		Users:   NewMockUsers(),
		Wallets: NewMockWallets(),
	}
}
