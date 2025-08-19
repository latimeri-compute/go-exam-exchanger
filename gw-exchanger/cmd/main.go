package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/server"
	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages/postgres"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type dbOptions struct {
	dbName     string
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     int
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	// err = godotenv.Load("../config.env")
	// if err != nil {
	// 	logger.Fatal(err.Error())
	// }

	var DBcfg dbOptions
	var serverCfg server.Config
	flag.IntVar(&serverCfg.Port, "SERVER_PORT", 8080, "порт сервера API")
	flag.StringVar(&serverCfg.Host, "SERVER_ADDRESS", os.Getenv("SERVER_ADDRESS"), "адрес сервера API")
	flag.StringVar(&DBcfg.dbUser, "POSTGRES_USER", os.Getenv("POSTGRES_USER"), "имя пользователя postgres")
	flag.StringVar(&DBcfg.dbPassword, "POSTGRES_PASSWORD", os.Getenv("POSTGRES_PASSWORD"), "пароль пользователя postgres")
	flag.StringVar(&DBcfg.dbName, "POSTGRES_DB", os.Getenv("POSTGRES_DB"), "название базы данных postgres")
	flag.StringVar(&DBcfg.dbHost, "POSTGRES_HOST", "db", "хост сервера postgres")
	flag.IntVar(&DBcfg.dbPort, "POSTGRES_PORT", 5432, "порт сервера postgres")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", DBcfg.dbUser, DBcfg.dbPassword, DBcfg.dbHost, DBcfg.dbPort, DBcfg.dbName)
	db, err := postgres.NewConnection(dsn, &gorm.Config{})
	if err != nil {
		sugar.Fatalf("Ошибка соединения с базой данных: %v", err)
	}
	logger.Info("соединение с базой данных установлено")

	srv := server.New(logger, db, serverCfg)

	sugar.Infof("Запуск сервера, порт: %d", serverCfg.Port)
	if err := srv.StartServer(); err != nil {
		sugar.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
