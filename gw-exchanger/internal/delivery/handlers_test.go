package delivery

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages/mocks"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// func init() {
// 	logger := zap.NewNop()
// 	h := NewHandler(logger, mocks.NewExchange())

// 	bufSize := 1024 * 1024
// 	lis := bufconn.Listen(bufSize)

// 	s := grpc.NewServer()
// 	pb.RegisterExchangeServiceServer(s, h)
// 	go func() {
// 		if err := s.Serve(lis); err != nil {
// 			log.Fatalf("Server exited with error: %v", err)
// 		}
// 	}()
// }

func TestGetExchangeRates(t *testing.T) {
	listener := NewTestListener()
	srv := NewTestServer()
	go func(t *testing.T) {
		if err := srv.Serve(listener); err != nil {
			t.Fatal(err)
		}
	}(t)
	defer srv.Stop()

	conn := NewTestConnection(t, listener)
	defer conn.Close()
	client := pb.NewExchangeServiceClient(conn)

	tests := []struct {
		name string
		want map[string]float32
	}{
		{
			want: map[string]float32{"rub->eur": 100, "rub->usd": 56, "usd->eur": 0.9},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			resp, err := client.GetExchangeRates(ctx, &pb.Empty{}, grpc.EmptyCallOption{})
			if err != nil {
				t.Fatal(err)
			}
			t.Log(reflect.TypeOf(resp))

			if !reflect.DeepEqual(resp.Rates, test.want) {
				t.Errorf("got: %v, want: %v", resp, test.want)
			}
		})
	}
}

func NewTestConnection(t *testing.T, listener *bufconn.Listener) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return listener.Dial()
		}),
	}

	conn, err := grpc.NewClient("passthrough:///bufnet", opts...)
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func NewTestListener() *bufconn.Listener {
	lis := bufconn.Listen(1024 * 1024)
	return lis
}

func NewTestServer() *grpc.Server {
	h := NewHandler(zap.NewNop(), mocks.NewExchange())
	s := grpc.NewServer()
	pb.RegisterExchangeServiceServer(s, h)
	return s
}
