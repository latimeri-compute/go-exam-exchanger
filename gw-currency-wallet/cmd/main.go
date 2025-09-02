package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/brocker"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/delivery"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/grpcclient"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/server"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages/postgres"
	"go.uber.org/zap"
)

//	@title			wallet API
//	@version		0.9
//	@description	wallet API supporting exchange between currencies

func main() {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	if err != nil {
		log.Fatal(err)
		return
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	err = godotenv.Load("../config.env")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	var DBcfg postgres.DBOptions
	serverConfig := server.Config{
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	postgres.InitFlags(&DBcfg)
	server.FlagInit(&serverConfig)

	gaddress := flag.String("gRPC address", "localhost:4000", "адрес удалённого grpc сервера")

	flag.Parse()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", DBcfg.DBUser, DBcfg.DBPassword, DBcfg.DBHost, DBcfg.DBPort, DBcfg.DBName)
	db, err := postgres.NewConnection(dsn)
	if err != nil {
		sugar.Fatalf("Ошибка соединения с базой данных: %v", err)
	}
	logger.Info("соединение с базой данных установлено")

	err = db.AutoMigrate(&storages.Wallet{}, &storages.User{})
	if err != nil {
		sugar.Errorf("Ошибка автомиграции: %v", err)
	}

	gclient, err := grpcclient.NewClient(*gaddress)
	if err != nil {
		sugar.Error("ошибка создания grpc клиента: ", err.Error())
		return
	}

	messageProducer, err := brocker.New(":9092")
	if err != nil {
		sugar.Error("Ошибка создания продюсера сообщений: ", err)
	}
	m := postgres.NewModels(db)
	h := delivery.NewHandler(m, logger.Sugar(), gclient, messageProducer, serverConfig.JWTSecret)

	srv := server.NewServer(h, logger.Sugar(), serverConfig)

	sugar.Infof("Запуск сервера, адрес: %s", srv.Server.Addr)
	err = srv.Serve()
	sugar.Error(err)
}
