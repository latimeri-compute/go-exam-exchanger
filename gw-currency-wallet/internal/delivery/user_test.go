package delivery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_storages "github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages/mocks"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name       string
		username   string
		email      string
		password   string
		want       string
		wantStatus int
	}{
		{
			name:       "верный пароль и email",
			email:      "new@email.com",
			password:   "password",
			want:       `{"message":"User registered successfully"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "повторяющийся email",
			email:      mock_storages.ValidUser.Email,
			password:   "password",
			want:       `{"error":"Username or email already exists"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "неправильно сформированнывй email",
			email:      "newemail",
			password:   "password",
			want:       `{"error":{"email":"is not a valid email"}}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	h := NewTestHandler("a", nil)
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

			status, body, _ := testutils.ReceiveResponse(t, client, http.MethodPost,
				"application/json", srv.URL+"/api/v1/register", sendJson)

			assert.Equal(t, test.want, string(body))
			assert.Equal(t, test.wantStatus, status)
		})
	}
}

func TestLoginUser(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		password     string
		wantContains string
		wantStatus   int
	}{
		{
			name:         "верный пользователь",
			email:        mock_storages.ValidUser.Email,
			password:     mock_storages.ValidPassword,
			wantContains: `"authentication_token"`,
			wantStatus:   http.StatusOK,
		},
		{
			name:         "неправильно сформированный email",
			email:        "newemail",
			password:     "password",
			wantContains: `{"error":{"email":"is not a valid email"}}`,
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "неверный пароль",
			email:        mock_storages.ValidUser.Email,
			password:     "wrongpassword",
			wantContains: `{"error":"Invalid username or password"}`,
			wantStatus:   http.StatusBadRequest,
		},
	}

	h := NewTestHandler("a", nil)
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

			status, body, _ := testutils.ReceiveResponse(t, client, http.MethodPost,
				"application/json", srv.URL+"/api/v1/login", sendJson)

			assert.Contains(t, string(body), test.wantContains)
			assert.Equal(t, test.wantStatus, status)
		})
	}
}
