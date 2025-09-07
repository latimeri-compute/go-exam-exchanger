package brocker

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
)

type TransactionMessage struct {
	WalletID     uint      `json:"wallet_id"`
	Type         string    `json:"type"`
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	AmountFrom   int       `json:"amount_from"`
	AmountTo     int       `json:"amount_to"`
	Timestamp    time.Time `json:"timestamp"`
}

type Producer struct {
	sarama.SyncProducer
}

func New(address string) (*Producer, error) {
	config := sarama.NewConfig()
	// config.Producer.Retry
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{address}, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer}, nil
}

func (p *Producer) MessageTransaction(transaction TransactionMessage) (partition int32, offset int64, err error) {
	v, err := json.Marshal(transaction)
	if err != nil {
		return 0, 0, err
	}
	msg := &sarama.ProducerMessage{
		Topic: "wallets_transactions",
		Value: sarama.StringEncoder(v),
	}

	return p.SendMessage(msg)
}
