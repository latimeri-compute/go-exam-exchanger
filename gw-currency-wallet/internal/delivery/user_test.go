package delivery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	mock_storages "github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages/mocks"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		want     string
	}{
		{
			name:     "valid",
			email:    "new@email.com",
			password: "password",
			want:     `{"message":"User registered successfully"}`,
		},
		{
			name:     "неверный email",
			email:    "newemail",
			password: "password",
			want:     `{"error":{"email":"is not a valid email"}}`,
		},
	}

	ctrl := gomock.NewController(t)
	mUsers := mock_storages.NewMockUserModelInterface(ctrl)
	mUsers.EXPECT().CreateUser(gomock.Any()).Return(nil).Times(1)

	mWallets := mock_storages.NewMockWalletModelInterface(ctrl)
	m := mock_storages.NewMockModels(mUsers, mWallets)
	defer ctrl.Finish()

	h := NewHandler(m, zap.NewNop().Sugar(), "a")
	srv := httptest.NewServer(Router(h))
	defer srv.Close()

	client := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sendJson, err := json.Marshal(loginJSON{
				Email:    test.email,
				Password: test.password,
			})
			if err != nil {
				t.Fatal(err)
			}

			body := testutils.ReceiveResponseBody(t, client, http.MethodPost,
				"application/json", srv.URL+"/api/v1/register", sendJson)

			assert.Equal(t, test.want, string(body))
		})
	}
}

// если бы мы знали что это такое но мы ведь не знаем что это такое
func TestLoginUser(t *testing.T) {
	// TODO
	tests := []struct {
		name     string
		email    string
		password string
		want     string
	}{
		{
			name:     "valid",
			email:    "new@email.com",
			password: "password",
			want:     `"error": "Invalid username or password"`,
		},
		{
			name:     "неверный email",
			email:    "newemail",
			password: "password",
			want:     `{"error":{"email":"is not a valid email"}}`,
		},
	}

	ctrl := gomock.NewController(t)
	mUsers := mock_storages.NewMockUserModelInterface(ctrl)
	mUsers.EXPECT().FindUser(&storages.User{Email: tests[0].email}).Return(nil).Times(1)
	// mUsers.EXPECT().FindUser(storages.User{Email: tests[1].email}).Return(storages.ErrRecordNotFound).Times(1)

	mWallets := mock_storages.NewMockWalletModelInterface(ctrl)
	m := mock_storages.NewMockModels(mUsers, mWallets)
	defer ctrl.Finish()

	h := NewHandler(m, zap.NewNop().Sugar(), "a")
	srv := httptest.NewServer(Router(h))
	defer srv.Close()

	client := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sendJson, err := json.Marshal(loginJSON{
				Email:    test.email,
				Password: test.password,
			})
			if err != nil {
				t.Fatal(err)
			}

			body := testutils.ReceiveResponseBody(t, client, http.MethodPost,
				"application/json", srv.URL+"/api/v1/login", sendJson)

			assert.Equal(t, test.want, string(body))
		})
	}
}
