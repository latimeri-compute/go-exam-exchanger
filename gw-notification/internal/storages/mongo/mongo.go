package mongo

import (
	"context"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/utils"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewWalletClient(client *mongo.Client) *WalletClient {
	coll := client.Database("wallets").Collection("wallets_transactions")
	return &WalletClient{coll}
}

func (c *WalletClient) InsertTransaction(transaction storages.Transaction, ctx context.Context) (any, error) {
	res, err := c.InsertOne(ctx, transaction, options.InsertOne())
	if err != nil {
		return "", err
	}
	return res.InsertedID, nil
}

func (c *WalletClient) Get(transaction storages.Transaction, ctx context.Context) ([]storages.Transaction, error) {
	filters := utils.StructToBson(transaction)
	cursor, err := c.Find(ctx, filters)
	if err != nil {
		return nil, err
	}
	var results []storages.Transaction
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
