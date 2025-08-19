package delivery

import (
	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages/mocks"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewTestHandler() *Handler {
	return NewHandler(zap.NewNop(), mocks.NewExchange())
}

func NewTestServer() *grpc.Server {
	h := NewTestHandler()
	s := grpc.NewServer()
	pb.RegisterExchangeServiceServer(s, h)
	return s
}
