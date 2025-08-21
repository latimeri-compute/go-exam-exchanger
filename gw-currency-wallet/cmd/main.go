package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/delivery"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/server"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages/postgres"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	err = godotenv.Load("../config.env")
	if err != nil {
		logger.Error(err.Error())
	}
	var DBcfg postgres.DBOptions
	flag.StringVar(&DBcfg.DBUser, "POSTGRES_USER", os.Getenv("POSTGRES_USER"), "имя пользователя postgres")
	flag.StringVar(&DBcfg.DBPassword, "POSTGRES_PASSWORD", os.Getenv("POSTGRES_PASSWORD"), "пароль пользователя postgres")
	flag.StringVar(&DBcfg.DBName, "POSTGRES_DB", os.Getenv("POSTGRES_DB"), "название базы данных postgres")
	flag.StringVar(&DBcfg.DBHost, "POSTGRES_HOST", "localhost", "хост сервера postgres")
	flag.IntVar(&DBcfg.DBPort, "POSTGRES_PORT", 5432, "порт сервера postgres")

	var serverConfig server.Config
	flag.IntVar(&serverConfig.Port, "SERVER_PORT", 4001, "порт сервера API")

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

	m := postgres.NewModels(db)
	h := delivery.NewHandler(m, logger.Sugar())
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", serverConfig.Addr, serverConfig.Port),
		Handler: delivery.Router(h),

		ErrorLog: slog.NewLogLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		}), slog.LevelError),

		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	sugar.Infof("Запуск сервера, адрес: %s", srv.Addr)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
