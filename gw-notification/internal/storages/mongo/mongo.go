package mongo

import (
	"context"

	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/storages"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type WalletClient struct {
	*mongo.Collection
}

func NewConnection(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewWalletClient(client *mongo.Client) *WalletClient {
	coll := client.Database("wallets").Collection("wallets_transactions")
	return &WalletClient{coll}
}

func (c *WalletClient) Insert(transaction storages.Transaction, ctx context.Context) error {
	_, err := c.InsertOne(ctx, transaction)
	return err
}

func (c *WalletClient) Get(transaction *storages.Transaction, ctx context.Context) ([]storages.Transaction, error) {
	cursor, err := c.Find(ctx, transaction)
	if err != nil {
		return nil, err
	}
	var results []storages.Transaction
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
