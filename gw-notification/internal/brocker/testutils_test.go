package brocker

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/storages"
	"go.uber.org/zap"
)

type transactionMessage struct {
	WalletID     uint      `json:"wallet_id"`
	Type         string    `json:"type"`
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	AmountFrom   int       `json:"amount_from"`
	AmountTo     int       `json:"amount_to"`
	Timestamp    time.Time `json:"timestamp"`
}

var uri = []string{"localhost:9092"}
var topics = []string{"wallets_transactions_test"}

func newTestConsumers(t *testing.T, ch chan storages.Transaction) *ConsumerGroup {
	conn, err := NewConnection(uri, zap.NewNop().Sugar())
	if err != nil {
		t.Fatal(err)
	}
	go conn.Consume(ch, topics, context.Background())
	return conn
}

func newProducer(t *testing.T) sarama.SyncProducer {
	config := sarama.NewConfig()
	// config.Producer.Retry
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(uri, config)
	if err != nil {
		t.Fatal(err)
	}
	return producer
}

func send(t *testing.T, p sarama.SyncProducer, transaction transactionMessage) {
	v, err := json.Marshal(transaction)
	if err != nil {
		t.Fatal(err)
	}
	msg := &sarama.ProducerMessage{
		Topic: topics[0],
		Value: sarama.StringEncoder(v),
	}
	_, _, err = p.SendMessage(msg)
	if err != nil {
		t.Fatal(err)
	}
}
