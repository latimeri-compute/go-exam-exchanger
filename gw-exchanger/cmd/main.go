package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/server"
	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages/postgres"
	"go.uber.org/zap"
)

type dbOptions struct {
	dbName     string
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     int
}

var (
	DBcfg     dbOptions
	serverCfg server.Config
	dsn       string
)

func init() {
	flag.IntVar(&serverCfg.Port, "SERVER_PORT", 443, "порт сервера API")
	flag.StringVar(&serverCfg.Host, "SERVER_ADDRESS", os.Getenv("SERVER_ADDRESS"), "адрес сервера API")
	flag.StringVar(&DBcfg.dbUser, "POSTGRES_USER", os.Getenv("POSTGRES_USER"), "имя пользователя postgres")
	flag.StringVar(&DBcfg.dbPassword, "POSTGRES_PASSWORD", os.Getenv("POSTGRES_PASSWORD"), "пароль пользователя postgres")
	flag.StringVar(&DBcfg.dbName, "POSTGRES_DB", os.Getenv("POSTGRES_DB"), "название базы данных postgres")
	flag.StringVar(&DBcfg.dbHost, "POSTGRES_HOST", os.Getenv("POSTGRES_HOST"), "хост сервера postgres")
	flag.IntVar(&DBcfg.dbPort, "POSTGRES_PORT", 5432, "порт сервера postgres")
	flag.Parse()

	dsn = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", DBcfg.dbUser, DBcfg.dbPassword, DBcfg.dbHost, DBcfg.dbPort, DBcfg.dbName)
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	db, err := postgres.NewConnection(dsn)
	if err != nil {
		sugar.Fatal("Ошибка соединения с базой данных: ", err)
	}
	logger.Info("соединение с базой данных установлено")

	srv := server.New(logger, db, serverCfg)

	if err := srv.StartServer(); err != nil {
		sugar.Fatalf("Ошибка запуска сервера: ", err)
	}
}
