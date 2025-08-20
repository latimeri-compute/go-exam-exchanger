package delivery

import (
	"testing"

	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages/mocks"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewTestHandler(t *testing.T) *Handler {
	t.Helper()
	return NewHandler(zap.NewNop(), mocks.NewExchange())
}

func NewTestServer(t *testing.T) *grpc.Server {
	t.Helper()
	h := NewTestHandler(t)
	s := grpc.NewServer()
	pb.RegisterExchangeServiceServer(s, h)
	return s
}
