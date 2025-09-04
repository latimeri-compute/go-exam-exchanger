package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/storages"
	"github.com/stretchr/testify/assert"
)

var testConn string = "mongodb://localhost:27017/?appName=mongodb-vscode+1.13.3&directConnection=true&serverSelectionTimeoutMS=2000"

func TestInsert(t *testing.T) {
	tests := []struct {
		name        string
		transaction storages.Transaction
	}{
		{
			name: "withdrawal",
			transaction: storages.Transaction{
				WalletID:     1,
				Type:         "withdrawal",
				FromCurrency: "usd",
				AmountFrom:   45000000,
				Timestamp:    time.Now(),
			},
		},
		{
			name: "exchanged between currencies",
			transaction: storages.Transaction{
				WalletID:     1,
				Type:         "exchange",
				FromCurrency: "usd",
				ToCurrency:   "rub",
				AmountFrom:   80000,
				AmountTo:     4670000000,
				Timestamp:    time.Now(),
			},
		},
	}

	m, err := newTestClient(t, testConn)
	if err != nil {
		t.Fatal(err)
	}
	defer m.Drop(context.Background())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			id, err := m.Insert(test.transaction, ctx)
			assert.NoError(t, err)
			assert.NotEmpty(t, id)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name        string
		transaction storages.Transaction
		want        any
	}{
		{
			name: "",
			transaction: storages.Transaction{
				FromCurrency: "usd",
			},
		},
		{
			transaction: storages.Transaction{
				Type:         "withdraw",
				FromCurrency: "rub",
			},
		},
	}

	m, err := newTestClient(t, testConn)
	if err != nil {
		t.Fatal(err)
	}
	fill(t, m)
	defer m.Drop(context.Background())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			res, err := m.Get(test.transaction, ctx)
			assert.NoError(t, err)
			t.Log(res)
		})
	}
}
