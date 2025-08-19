package delivery

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/pkg/testutils"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"gorm.io/gorm"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetExchangeRates(t *testing.T) {
	listener := testutils.NewTestListener()
	srv := NewTestServer()
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
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(resp.Rates, test.want) {
				t.Errorf("got: %v, want: %v", resp.Rates, test.want)
			}
		})
	}
}

func TestGetExchangeRateForCurrency(t *testing.T) {
	listener := testutils.NewTestListener()
	srv := NewTestServer()
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
			wantError:    status.Error(codes.NotFound, gorm.ErrRecordNotFound.Error()),
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
			if err != nil && !errors.Is(err, test.wantError) {
				t.Errorf("got error: %v, want: %v", err, test.wantError)
			}

			// ew ugly
			w := struct {
				ToCurrency   string
				FromCurrency string
				Rate         float32
			}{
				ToCurrency:   test.toCurrency,
				FromCurrency: test.fromCurrency,
				Rate:         test.rate,
			}

			if test.wantError == nil && (resp.FromCurrency != test.fromCurrency ||
				resp.ToCurrency != test.toCurrency ||
				resp.Rate != test.rate) {
				t.Errorf("got: %v, want: %v", resp, w)
			}
		})
	}
}
