package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "gitgub.com/latimeri-compute/gw-exchanger/internal/delivery/exchange"
	delivery "gitgub.com/latimeri-compute/gw-exchanger/internal/delivery/response"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ServerGRPC struct {
	Logger *zap.Logger
	Server *grpc.Server
	Cfg    Config
}

type Config struct {
	Port int
}

func New(logger *zap.Logger, port int) *ServerGRPC {
	srv := grpc.NewServer()
	ex := delivery.NewHandler(logger)
	pb.RegisterExchangeServiceServer(srv, ex)

	return &ServerGRPC{
		Logger: logger,
		Server: srv,
	}
}

func (srv *ServerGRPC) StartServer() error {
	go func() {
		// create a quit channel
		quit := make(chan os.Signal, 1)

		// listen for incoming SIGINT and SIGTERM
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// read signal from the quit channel. Will block until signal is received
		signal := <-quit

		srv.Logger.Sugar().Info("shutting down server", "signal", signal.String())

		srv.Server.GracefulStop()

	}()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.Cfg.Port))
	if err != nil {
		return err
	}
	defer listen.Close()

	srv.Logger.Info("Запуск сервера")
	err = srv.Server.Serve(listen)
	if !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}

	return nil
}

// 	go func() {
// 		// create a quit channel
// 		quit := make(chan os.Signal, 1)

// 		// listen for incoming SIGINT and SIGTERM
// 		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

// 		// read signal from the quit channel. Will block until signal is received
// 		s := <-quit

// 		app.logger.Info("shutting down server", "signal", s.String())

// 		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 		defer cancel()

// 		err := srv.Shutdown(ctx)
// 		if err != nil {
// 			shutdownError <- err
// 		}

// 		app.logger.Info("completing background tasks", "addr", srv.Addr)

// 		app.waitgroup.Wait()
// 		shutdownError <- nil
// 	}()

// 	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

// 	err := srv.ListenAndServe()
// 	if !errors.Is(err, http.ErrServerClosed) {
// 		return err
// 	}

// 	err = <-shutdownError
// 	if err != nil {
// 		return err
// 	}

// 	app.logger.Info("stopped server", "addr", srv.Addr)

// 	return nil
// }
