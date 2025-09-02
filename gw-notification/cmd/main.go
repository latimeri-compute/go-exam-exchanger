package cmd

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/brocker"
	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/storages/mongo"
	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/server"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	var producerAddr string
	var mongoUri string
	flag.StringVar(&producerAddr, "адрес продюсера сообщений", "", "")
	flag.StringVar(&mongoUri, "строка подключения MongoDB", "", "")
	flag.Parse()

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
