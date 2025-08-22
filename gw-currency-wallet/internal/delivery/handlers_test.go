package delivery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_storages "github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages/mocks"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

// func TestMain(m *testing.M) {
// 	ctrl := gomock.NewController(m)
// 	mod := mock_storages.NewMockModels(ctrl)
// 	h := NewHandler(mod, zap.NewNop().Sugar())

// 	srv := httptest.NewServer(Router(h))
// 	defer srv.Close()

// 	code := m.Run()

// 	os.Exit(code)
// }

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
	}

	ctrl := gomock.NewController(t)
	mUsers := mock_storages.NewMockUserModelInterface(ctrl)
	mWallets := mock_storages.NewMockWalletModelInterface(ctrl)
	m := mock_storages.NewMockModels(mUsers, mWallets)
	defer ctrl.Finish()

	h := NewHandler(m, zap.NewNop().Sugar(), "a")
	srv := httptest.NewServer(Router(h))
	defer srv.Close()

	client := &http.Client{}

	mUsers.EXPECT().CreateUser(gomock.Any()).Return(nil).Times(1)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sendJson, err := json.Marshal(loginJSON{
				Email:    test.email,
				Password: test.password,
			})
			if err != nil {
				t.Fatal(err)
			}

			body := testutils.ReceiveResponseBody(t, client, http.MethodPost, srv.URL+"/api/v1/register", sendJson)
			t.Log(string(body))
			assert.Equal(t, test.want, string(body))
		})
	}

}
