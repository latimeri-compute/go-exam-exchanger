package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"gitgub.com/latimeri-compute/gw-exchanger/internal/server"
	"go.uber.org/zap"
)

type options struct {
	port int
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

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.port))
	if err != nil {
		sugar.Fatalf("Ошибка прослушивания порта %d: %s", cfg.port, err)
	}
	defer listen.Close()

	srv := server.New(logger, cfg.port)

	sugar.Infof("Запуск сервера", "порт", cfg.port)
	if err := srv.StartServer(); err != nil {
		os.Exit(1)
		sugar.Fatal("Ошибка запуска сервера:", err)
	}

}
