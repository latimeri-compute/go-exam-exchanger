package brocker

import (
	"encoding/json"
	"errors"

	"github.com/IBM/sarama"
	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/storages"
	"go.uber.org/zap"
)

var ErrClosed = errors.New("Channel closed")
var ErrExpiredContext = errors.New("Context expired")

type consumer struct {
	transactionsExchange chan storages.Transaction
	errorsChannel        chan error
	logger               *zap.SugaredLogger
	transactionDB        storages.WalletModelInterface
}

func newConsumer(logger *zap.SugaredLogger, transactionsExchange chan storages.Transaction) *consumer {
	return &consumer{
		transactionsExchange: transactionsExchange,
		logger:               logger,
	}
}

func (c *consumer) Setup(sarama.ConsumerGroupSession) error {

	//TODO
	c.logger.Info("Message consumers set up!")
	return nil
}

func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	//TODO
	c.logger.Info("Message consumers cleaned!")
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	//TODO
	for {
		select {
		case mes, ok := <-claim.Messages():
			if !ok {
				c.logger.Info("Канал сообщений консьюмера закрыт")
				return nil
			}

			c.logger.Info("получено сообщение: ", &mes)

			var transaction *storages.Transaction
			err := json.Unmarshal(mes.Value, &transaction)
			if err != nil {
				c.logger.Error("Ошибка анмаршиллинга: ", err)
				continue
			}

			c.transactionsExchange <- *transaction
			// стоит ли помечать в случае ошибки добавления в бд?..
			session.MarkMessage(mes, "")

		case <-session.Context().Done():
			return ErrExpiredContext
		}
	}
}
