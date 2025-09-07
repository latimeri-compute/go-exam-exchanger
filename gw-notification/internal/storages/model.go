package storages

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("Document not found")

type WalletModelInterface interface {
	InsertTransaction(transaction Transaction, ctx context.Context) (any, error)
	Get(transaction Transaction, ctx context.Context) ([]Transaction, error)
}

type Transaction struct {
	WalletID     uint      `bson:"wallet_id" json:"wallet_id"`
	Type         string    `bson:"type" json:"type"`
	AmountFrom   int       `bson:"amount_from" json:"amount_from"`
	AmountTo     int       `bson:"amount_to,omitempty" json:"amount_to,omitempty"`
	FromCurrency string    `bson:"from_currency" json:"from_currency"`
	ToCurrency   string    `bson:"to_currency,omitempty" json:"to_currency"`
	Timestamp    time.Time `bson:"timestamp" json:"timestamp"`
}
