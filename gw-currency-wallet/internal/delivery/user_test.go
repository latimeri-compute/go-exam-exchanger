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
			username:   "blahblah",
			password:   "password",
			want:       `{"message":"User registered successfully"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "повторяющийся email",
			email:      mock_storages.ValidUser.Email,
			password:   "password",
			username:   "blahblah",
			want:       `{"error":"Username or email already exists"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "повторяющийся username",
			username:   mock_storages.ValidUser.Username,
			password:   "password",
			email:      "new@email.com",
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

	h := NewTestHandler("a", nil, nil)
	srv := httptest.NewServer(Router(h))
	defer srv.Close()

	client := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sendJson, err := json.Marshal(registerJSON{
				Email:    test.email,
				Username: test.username,
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
		username     string
		password     string
		wantContains string
		wantStatus   int
	}{
		{
			name:         "верный пользователь",
			username:     mock_storages.ValidUser.Username,
			password:     mock_storages.ValidPassword,
			wantContains: `"authentication_token"`,
			wantStatus:   http.StatusOK,
		},
		{
			name:         "неверный пароль",
			username:     mock_storages.ValidUser.Username,
			password:     "wrongpassword",
			wantContains: `{"error":"Invalid username or password"}`,
			wantStatus:   http.StatusUnauthorized,
		},
	}

	h := NewTestHandler("a", nil, nil)
	srv := httptest.NewServer(Router(h))
	defer srv.Close()

	client := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sendJson, err := json.Marshal(loginJSON{
				Username: test.username,
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
