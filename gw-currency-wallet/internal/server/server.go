package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/delivery"
	"go.uber.org/zap"
)

type Config struct {
	Port         int
	Addr         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	JWTSecret string
}

type Server struct {
	Cfg       Config
	Server    *http.Server
	logger    *zap.SugaredLogger
	waitgroup sync.WaitGroup
}

func NewServer(h *delivery.Handler, logger *zap.SugaredLogger, cfg Config) *Server {

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		Handler:      delivery.Router(h),
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,

		ErrorLog: slog.NewLogLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		}), slog.LevelError),
	}
	Server := &Server{
		Cfg:    cfg,
		Server: srv,
		logger: logger,
	}
	return Server
}

// TODO graceful stopping?

func (s *Server) Serve() error {

	shutdownError := make(chan error)

	go func() {
		// create a quit channel
		quit := make(chan os.Signal, 1)

		// listen for incoming SIGINT and SIGTERM
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// read signal from the quit channel. Will block until signal is received
		signal := <-quit

		s.logger.Info("shutting down server", "signal", signal.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := s.Server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		s.logger.Info("completing background tasks", "addr", s.Server.Addr)

		s.waitgroup.Wait()
		shutdownError <- nil
	}()

	s.logger.Info("starting server", "addr", s.Cfg.Addr, "port", s.Cfg.Port)

	err := s.Server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	s.logger.Info("stopped server", "addr", s.Server.Addr)

	return nil
}

// отправить функцию на параллель
func (s *Server) Background(fn func()) {
	s.waitgroup.Add(1)

	go func() {
		defer s.waitgroup.Done()

		defer func() {
			if err := recover(); err != nil {
				s.logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		fn()
	}()
}
