package grpcclient

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

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

// func NewServer() {
// 	h := NewHandler()
// 	srv := grpc.NewServer()
// }

func NewTestListener() *bufconn.Listener {
	lis := bufconn.Listen(1024 * 1024)
	return lis
}
