package main

import (
	"context"
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
	producerAddr []string = []string{"kafka:9093"}
	mongoUri     string   = "mongodb://mongodb/wallets"
	// dbUser       string
	// dbPassword   string
	// dbPort       int
)

// func init() {
// 	// flag.StringVar(&dbUser, "MONGO_USER", os.Getenv("MONGO_USER"), "пользователь MongoDB")
// 	// flag.StringVar(&dbPassword, "MONGO_PASSWORD", os.Getenv("MONGO_PASSWORD"), "пароль MongoDB")
// 	flag.IntVar(&dbPort, "MONGO_PORT", 27017, "порт MongoDB")
// 	flag.Parse()

// 	// mongoUri = fmt.Sprintf("mongodb://mongodb/", dbPort)
// }

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	sarama.Logger = log.New(os.Stdout, "[sarama]", log.LstdFlags)

	consumer, err := brocker.NewConnection(producerAddr, sugar)
	if err != nil {
		sugar.Error("Ошибка соединения с брокером: ", err)
		return
	}
	defer consumer.Group.Close()
	sugar.Info("соединение с брокером установлено")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	sugar.Debug("Mongo URI: ", mongoUri)
	mdb, err := mongo.NewConnection(mongoUri)
	if err != nil {
		sugar.Error("Ошибка соединения с базой данных: ", err)
		return
	}
	defer mdb.Disconnect(ctx)
	sugar.Info("Соединение с базой данных установлено")
	wallets := mongo.NewWalletClient(mdb)

	sugar.Info("запуск сервера...")
	srv := server.New(consumer, wallets, sugar)
	srv.Start()
}
