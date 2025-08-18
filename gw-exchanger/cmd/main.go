package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/server"
	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages/postgres"
	"go.uber.org/zap"
)

type options struct {
	port int
	dsn  string
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	logger.Info("logger initialized")

	var cfg options
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.dsn, "dsn", "DSN string", "")

	db, err := postgres.NewConnection(cfg.dsn, nil)
	if err != nil {
		sugar.Fatalf("Ошибка соединения с базой данных: %v", err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.port))
	if err != nil {
		sugar.Fatalf("Ошибка прослушивания порта %d: %s", cfg.port, err)
	}
	defer listen.Close()

	srv := server.New(logger, cfg.port, db)

	sugar.Infof("Запуск сервера", "порт", cfg.port)
	if err := srv.StartServer(); err != nil {
		os.Exit(1)
		sugar.Fatal("Ошибка запуска сервера:", err)
	}

}
