package server

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/brocker"
	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/storages"
	"go.uber.org/zap"
)

type Server struct {
	Consumer     *brocker.ConsumerGroup
	Transactions storages.WalletModelInterface
	logger       *zap.SugaredLogger
	wg           sync.WaitGroup
}

func New(consumer *brocker.ConsumerGroup, transactions storages.WalletModelInterface) *Server {
	return &Server{
		Consumer:     consumer,
		Transactions: transactions,
	}
}

func (s *Server) Start() {
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		signal := <-quit

		s.logger.Infof("Получен сигнал: %s; Остановка приложения...", signal.String())

		s.wg.Wait()
		err := s.Consumer.Group.Close()
		if err != nil {
			shutdownError <- err
		}
		shutdownError <- nil
	}()

	exchangeCh := make(chan storages.Transaction)

	go s.Consumer.Consume(exchangeCh, []string{}, context.Background())

	for transaction := range exchangeCh {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()

			err := s.Transactions.Insert(transaction, ctx)
			if err != nil {
				s.logger.DPanic("Ошибка добавления документа в базу данных: ", err)
				s.logger.Error("Ошибка добавления документа в базу данных: ", err)
			}
		}()
	}

	if err := <-shutdownError; err != nil {
		s.logger.Error("Ошибка остановки приложения: ", err.Error())
	}
}
