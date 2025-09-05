package delivery

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/pkg/testutils"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type receivedExchange struct {
	ToValuteCode   string
	FromValuteCode string
	Rate           float32
}

func TestGetExchangeRates(t *testing.T) {
	listener := testutils.NewTestListener()
	srv := NewTestServer(t)
	go func(t *testing.T) {
		if err := srv.Serve(listener); err != nil {
			t.Error(err)
		}
	}(t)
	defer srv.Stop()

	conn := testutils.NewTestConnection(t, listener)
	defer conn.Close()
	client := pb.NewExchangeServiceClient(conn)

	tests := []struct {
		want map[string]float32
	}{
		{
			want: map[string]float32{"rub->eur": 100, "rub->usd": 56, "usd->eur": 0.9, "usd->rub": 0.0037},
		},
	}

	for ti, test := range tests {
		t.Run(fmt.Sprintf("%02d", ti), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			resp, err := client.GetExchangeRates(ctx, &pb.Empty{}, grpc.EmptyCallOption{})
			assert.NoError(t, err)
			assert.Equal(t, resp.Rates, test.want)
		})
	}
}

func TestGetExchangeRateForCurrency(t *testing.T) {
	listener := testutils.NewTestListener()
	srv := NewTestServer(t)
	go func(t *testing.T) {
		if err := srv.Serve(listener); err != nil {
			t.Error(err)
		}
	}(t)
	defer srv.Stop()

	conn := testutils.NewTestConnection(t, listener)
	defer conn.Close()
	client := pb.NewExchangeServiceClient(conn)

	tests := []struct {
		name         string
		fromCurrency string
		toCurrency   string
		rate         float32
		wantError    error
	}{
		{
			name:         "существующие валюты",
			fromCurrency: "usd",
			toCurrency:   "rub",
			rate:         0.0037,
		},
		{
			name:         "несуществующие валюты",
			fromCurrency: "sda",
			toCurrency:   "what",
			rate:         00,
			wantError:    status.Error(codes.NotFound, storages.ErrNotFound.Error()),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			mes := &pb.CurrencyRequest{
				FromCurrency: test.fromCurrency,
				ToCurrency:   test.toCurrency,
			}
			resp, err := client.GetExchangeRateForCurrency(ctx, mes, grpc.EmptyCallOption{})
			if test.wantError != nil {
				assert.ErrorIs(t, err, test.wantError)
			} else {
				assert.Equal(t, test.fromCurrency, resp.FromCurrency)
				assert.Equal(t, test.toCurrency, resp.ToCurrency)
				assert.Equal(t, test.rate, resp.Rate)
			}
		})
	}
}
