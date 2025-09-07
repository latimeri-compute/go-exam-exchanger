package postgres

import (
	"testing"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DSN string = "postgresql://postgres:password@localhost:5432/wallets_test"

func setupDB(t *testing.T, db *gorm.DB) {
	err := db.AutoMigrate(storages.User{}, storages.Wallet{})
	if err != nil {
		t.Fatal(err)
	}
}

func teardownDB(t *testing.T, db *gorm.DB) {
	err := db.Exec("DROP TABLE users; DROP TABLE wallets;").Error
	if err != nil {
		t.Fatal(err)
	}
}

func newTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Discard,
	})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func createOneWallet(t *testing.T, db *gorm.DB) storages.Wallet {
	var wallet storages.Wallet
	err := db.Create(&wallet).Error
	if err != nil {
		t.Fatal(err)
	}

	return wallet
}
