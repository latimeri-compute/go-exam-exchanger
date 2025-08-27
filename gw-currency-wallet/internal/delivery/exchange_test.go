package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/grpcclient"
	mock_storages "github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages/mocks"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/testutils"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TestGetRates(t *testing.T) {
	tests := []struct {
		name       string
		userId     int
		wantBody   string
		wantStatus int
	}{
		{
			name:       "существующий пользователь",
			userId:     1,
			wantBody:   `{"rates":`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "несуществующий пользователь",
			userId:     99,
			wantBody:   `{"error":"Internal server error"}`,
			wantStatus: http.StatusInternalServerError,
		},
	}

	lis := grpcclient.NewTestListener()
	grpcSrv := grpc.NewServer()
	pb.RegisterExchangeServiceServer(grpcSrv, grpcclient.NewHandler())
	go func(t *testing.T) {
		if err := grpcSrv.Serve(lis); err != nil {
			t.Error(err)
		}
	}(t)
	defer grpcSrv.Stop()

	conn := grpcclient.NewTestConnection(t, lis)
	defer conn.Close()

	jwtString := "string!"
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	m := mock_storages.NewMockModels()
	h := NewHandler(m, logger.Sugar(), pb.NewExchangeServiceClient(conn), jwtString)
	// h := NewTestHandler(jwtString, pb.NewExchangeServiceClient(conn))
	srv := httptest.NewServer(Router(h))
	defer srv.Close()

	client := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status, body, _ := testutils.RequestWithJWT(
				t, client,
				[]byte(jwtString), []byte{},
				http.MethodGet, "application/json", srv.URL+"/api/v1/exchange/rates",
				test.userId)
			assert.Contains(t, string(body), test.wantBody)
			assert.Equal(t, test.wantStatus, status)
		})
	}
}

func TestExchangeFunds(t *testing.T) {
	tests := []struct {
		name       string
		userId     int
		from       string
		to         string
		amount     float64
		wantBody   string
		wantStatus int
	}{
		{
			name:       "существующий пользователь",
			userId:     1,
			from:       "usd",
			to:         "rub",
			amount:     20,
			wantBody:   `{"message":"Exchange successful","exchanged_amount":1700,"new_balance":`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "несуществующий пользователь",
			userId:     99,
			from:       "usd",
			to:         "rub",
			amount:     20,
			wantBody:   `{"error":"Internal server error"}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "сумма превышает баланс",
			userId:     1,
			from:       "usd",
			to:         "rub",
			amount:     20000,
			wantBody:   `{"error":"Insufficient funds or invalid currencies"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	lis := grpcclient.NewTestListener()
	grpcSrv := grpc.NewServer()
	pb.RegisterExchangeServiceServer(grpcSrv, grpcclient.NewHandler())
	go func(t *testing.T) {
		if err := grpcSrv.Serve(lis); err != nil {
			t.Error(err)
		}
	}(t)
	defer grpcSrv.Stop()

	conn := grpcclient.NewTestConnection(t, lis)
	defer conn.Close()

	jwtString := "string!"
	h := NewTestHandler(jwtString, pb.NewExchangeServiceClient(conn))
	srv := httptest.NewServer(Router(h))
	defer srv.Close()

	client := &http.Client{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			json := testutils.JsonShortcut(t, exchangeRequest{
				FromCurrency: test.from,
				ToCurrency:   test.to,
				Amount:       test.amount,
			})

			status, body, _ := testutils.RequestWithJWT(
				t, client,
				[]byte(jwtString), json,
				http.MethodPost, "application/json", srv.URL+"/api/v1/exchange",
				test.userId)
			assert.Contains(t, string(body), test.wantBody)
			assert.Equal(t, test.wantStatus, status)
		})
	}
}
