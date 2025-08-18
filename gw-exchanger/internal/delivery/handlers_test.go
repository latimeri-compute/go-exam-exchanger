package delivery

import (
	"context"
	"testing"

	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages/mocks"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestGetExchangeRates(t *testing.T) {
	tests := []struct {
		name string
	}{
		{},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			// TODO
		})
	}
}

func NewTestServer() *grpc.Server {
	// TODO
	bufSize := 1024 * 1024
	var logger *zap.Logger
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	h := NewHandler(logger, &mocks.MockExchange{})
	pb.RegisterExchangeServiceServer(s, h)

	return s
}
