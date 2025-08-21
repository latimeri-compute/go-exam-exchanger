package mock_storages

import (
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
)

func NewMockModels(users *MockUserModelInterface, wallets *MockWalletModelInterface) *storages.Models {
	return &storages.Models{
		Users:   users,
		Wallets: wallets,
	}
}
