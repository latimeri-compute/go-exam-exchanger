package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	tests := []struct {
		name       string
		id         uint
		want       string
		wantStatus int
	}{
		{
			name:       "существующий пользователь",
			id:         1,
			want:       `{"balance":{"USD":0,"EUR":0,"RUB":0}}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "несуществующий пользователь",
			id:         99,
			want:       `{"error":"Internal server error"}`,
			wantStatus: http.StatusInternalServerError,
		},
	}

	jwtString := "string!"
	h := NewTestHandler(jwtString, nil)
	srv := httptest.NewServer(Router(h))
	defer srv.Close()

	client := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status, body, _ := testutils.RequestWithJWT(
				t, client,
				[]byte(jwtString), nil,
				http.MethodGet, "application/json", srv.URL+"/api/v1/balance",
				int(test.id))
			assert.Equal(t, test.want, string(body))
			assert.Equal(t, test.wantStatus, status)
		})
	}

}

func TestDeposit(t *testing.T) {
	tests := []struct {
		name       string
		userId     int
		want       string
		amount     float64
		currency   string
		wantStatus int
	}{
		{
			name:       "существующий пользователь",
			userId:     1,
			amount:     80.97,
			currency:   "rub",
			want:       `{"message":"deposit successful","new_balance":{"USD":100,"EUR":100,"RUB":180.97}}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "несуществующий пользователь",
			userId:     99,
			amount:     100,
			currency:   "usd",
			want:       `{"error":"Internal server error"}`,
			wantStatus: http.StatusInternalServerError,
		},
	}

	jwtString := "string!"
	h := NewTestHandler(jwtString, nil)
	srv := httptest.NewServer(Router(h))
	defer srv.Close()

	client := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			json := testutils.JsonShortcut(t, fundsRequest{
				Amount:   test.amount,
				Currency: test.currency,
			})

			status, body, _ := testutils.RequestWithJWT(
				t, client,
				[]byte(jwtString), json,
				http.MethodPost, "application/json", srv.URL+"/api/v1/deposit",
				test.userId)
			assert.Equal(t, test.want, string(body))
			assert.Equal(t, test.wantStatus, status)
		})
	}
}

func TestWithdraw(t *testing.T) {
	tests := []struct {
		name       string
		userId     int
		want       string
		amount     float64
		currency   string
		wantStatus int
	}{
		{
			name:       "существующий пользователь",
			userId:     1,
			amount:     80.97,
			currency:   "rub",
			want:       `{"message":"withdrawal successful","new_balance":{"USD":100,"EUR":100,"RUB":19.03}}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "несуществующий пользователь",
			userId:     99,
			amount:     100,
			currency:   "usd",
			want:       `{"error":"Internal server error"}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "слишком большая сумма для снятия",
			userId:     1,
			amount:     1000000,
			currency:   "usd",
			want:       `{"error":"Insufficient funds or invalid amount"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	jwtString := "string!"
	h := NewTestHandler(jwtString, nil)
	srv := httptest.NewServer(Router(h))
	defer srv.Close()

	client := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			json := testutils.JsonShortcut(t, fundsRequest{
				Amount:   test.amount,
				Currency: test.currency,
			})

			status, body, _ := testutils.RequestWithJWT(
				t, client,
				[]byte(jwtString), json,
				http.MethodPost, "application/json", srv.URL+"/api/v1/withdraw",
				test.userId)
			assert.Equal(t, test.want, string(body))
			assert.Equal(t, test.wantStatus, status)
		})
	}
}
