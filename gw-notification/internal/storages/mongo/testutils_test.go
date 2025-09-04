package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/storages"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func newTestClient(t *testing.T, uri string) (*WalletClient, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}

	coll := client.Database("wallets_test").Collection("wallets_transactions")
	return &WalletClient{coll}, nil
}

func fill(t *testing.T, c *WalletClient) {
	var transactions = []storages.Transaction{
		{
			WalletID:     2,
			Type:         "exchange",
			FromCurrency: "rub",
			ToCurrency:   "eur",
			AmountFrom:   90212311,
			AmountTo:     12111222000000,
			Timestamp:    time.Now(),
		},
		{
			WalletID:     2,
			Type:         "exchange",
			FromCurrency: "rub",
			ToCurrency:   "eur",
			AmountFrom:   90212311,
			AmountTo:     12111222000000,
			Timestamp:    time.Now(),
		},
		{
			WalletID:     1,
			Type:         "withdraw",
			FromCurrency: "rub",
			AmountFrom:   90212311,
			Timestamp:    time.Now(),
		},
		{
			WalletID:     7,
			Type:         "deposit",
			FromCurrency: "usd",
			AmountFrom:   900000,
			Timestamp:    time.Now(),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, err := c.InsertMany(ctx, transactions)
	if err != nil {
		t.Fatal(err)
	}
}
