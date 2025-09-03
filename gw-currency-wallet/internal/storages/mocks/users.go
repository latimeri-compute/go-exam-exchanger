package mock_storages

import (
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ValidPassword string         = "passwordhehehehhe"
	ValidUser     *storages.User = &storages.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:    "admin@admin.admin",
		WalletID: 1,
		Wallet:   ValidWallet,
	}
)

type MockUsers struct{}

func (m *MockUsers) CreateUser(user *storages.User) error {
	if user.Email == ValidUser.Email {
		return storages.ErrRecordExists
	}
	return nil
}
func (m *MockUsers) FindUser(user *storages.User) error {
	defaultPass, err := bcrypt.GenerateFromPassword([]byte(ValidPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	if user.Email == ValidUser.Email || user.ID == ValidUser.ID {
		user.ID = ValidUser.ID
		user.Email = ValidUser.Email
		user.Model = ValidUser.Model
		user.Wallet = ValidUser.Wallet
		user.Wallet.ID = ValidUser.Wallet.ID
		user.WalletID = ValidUser.WalletID
		user.PasswordHash = defaultPass
	} else {
		return storages.ErrRecordNotFound
	}
	return nil
}

func NewMockUsers() *MockUsers {
	return &MockUsers{}
}
