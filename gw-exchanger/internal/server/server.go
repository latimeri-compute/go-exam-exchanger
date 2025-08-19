package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	delivery "github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/delivery"
	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ServerGRPC struct {
	Logger *zap.Logger
	Server *grpc.Server
	cfg    Config
}

type Config struct {
	Port int
	Host string
}

func New(logger *zap.Logger, db storages.ExchangerModelInterface, cfg Config) *ServerGRPC {
	srv := grpc.NewServer()
	ex := delivery.NewHandler(logger, db)
	pb.RegisterExchangeServiceServer(srv, ex)

	return &ServerGRPC{
		Logger: logger,
		Server: srv,
		cfg:    cfg,
	}
}

func (srv *ServerGRPC) StartServer() error {
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		signal := <-quit

		srv.Logger.Sugar().Info("остановка сервера", "сигнал", signal.String())
		srv.Server.GracefulStop()
	}()

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", srv.cfg.Host, srv.cfg.Port))
	if err != nil {
		return err
	}
	defer listen.Close()

	err = srv.Server.Serve(listen)
	if !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}

	return nil
}
