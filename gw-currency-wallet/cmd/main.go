package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/brocker"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/delivery"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/grpcclient"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/server"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages/postgres"
	proto_exchange "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

//	@title			wallet API
//	@version		0.9
//	@description	wallet API supporting exchange between currencies

var (
	dbCfg        postgres.DBOptions
	serverConfig = server.Config{
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	dsn          = ""
	grpcAddress  = ""
	kafkaAddress = "gw-notification_default:4010"
)

func init() {
	postgres.InitFlags(&dbCfg)
	server.FlagInit(&serverConfig)
	flag.StringVar(&grpcAddress, "gRPC address", "gw-exchanger", "адрес удалённого grpc сервера")
	flag.Parse()
	dsn = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", dbCfg.DBUser, dbCfg.DBPassword, dbCfg.DBHost, dbCfg.DBPort, dbCfg.DBName)
}

func main() {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	if err != nil {
		log.Fatal(err)
		return
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Debug("grpc address: ", grpcAddress)
	sugar.Debug("dsn: ", dsn)

	db, err := postgres.NewConnection(dsn)
	if err != nil {
		sugar.Fatal("Ошибка соединения с базой данных: ", err)
	}
	logger.Info("соединение с базой данных установлено")

	err = db.AutoMigrate(&storages.Wallet{}, &storages.User{})
	if err != nil {
		sugar.Error("Ошибка автомиграции: ", err)
	}

	gclient, err := grpcclient.NewClient(grpcAddress)
	if err != nil {
		sugar.Error("ошибка создания grpc клиента: ", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = gclient.GetExchangeRates(ctx, &proto_exchange.Empty{}, grpc.EmptyCallOption{})
	if err != nil {
		sugar.Error("ошибка соединения с удалённым сервером обмена валют: ", err)
	}

	messageProducer, err := brocker.New("kafka:9093")
	if err != nil {
		sugar.Error("Ошибка создания продюсера сообщений: ", err)
	}

	m := postgres.NewModels(db)
	h := delivery.NewHandler(m, logger.Sugar(), gclient, messageProducer, serverConfig.JWTSecret)

	srv := server.NewServer(h, logger.Sugar(), serverConfig)
	sugar.Info("Запуск сервера, адрес: ", srv.Server.Addr)
	if err = srv.Serve(); err != nil {
		sugar.Error(err)
	}
}
