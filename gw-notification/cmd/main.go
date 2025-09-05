package cmd

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/brocker"
	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/server"
	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/storages/mongo"
	"go.uber.org/zap"
)

var (
	producerAddr string
	mongoUri     string
)

func init() {
	flag.StringVar(&producerAddr, "адрес продюсера сообщений", "", "")
	flag.StringVar(&mongoUri, "строка подключения MongoDB", "", "")
	flag.Parse()
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sarama.Logger = log.New(os.Stdout, "[sarama]", log.LstdFlags)

	consumer, err := brocker.NewConnection([]string{producerAddr})
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer consumer.Group.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	mdb, err := mongo.NewConnection(mongoUri)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer mdb.Disconnect(ctx)

	wallets := mongo.NewWalletClient(mdb)

	srv := server.New(consumer, wallets)
	srv.Start()
}
