package brocker

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/storages"
	"go.uber.org/zap"
)

type ConsumerGroup struct {
	Group  sarama.ConsumerGroup
	logger *zap.SugaredLogger
}

func NewConnection(addr []string) (*ConsumerGroup, error) {

	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}

	con, err := sarama.NewConsumerGroup(addr, "wallet_consume", cfg)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{Group: con}, err
}

func (c *ConsumerGroup) Consume(receiver chan storages.Transaction, topics []string, ctx context.Context) {
	consumer := newConsumer(c.logger, receiver)
	for {
		// TODO написать клэймы
		if err := c.Group.Consume(ctx, topics, consumer); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				c.logger.Info("Consumer group closed")
				return
			} else {
				c.logger.DPanic(err)
				c.logger.Error(err)
			}
		}
		if ctx.Err() != nil {
			c.logger.Info("Consumer context expired")
			return
		}
	}
}
