package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/delivery/middleware"
	"golang.org/x/crypto/bcrypt"
)

func HashString(t *testing.T, str string) []byte {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	return bytes
}

func JsonShortcut(t *testing.T, mes any) []byte {
	json, err := json.Marshal(mes)
	if err != nil {
		t.Fatal(err)
	}
	return json
}

func ReceiveBodyPost(t *testing.T, urlPath string, body []byte) []byte {
	t.Helper()
	resp, err := http.Post(urlPath, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	b = bytes.TrimSpace(b)
	return b
}

func ReceiveResponse(t *testing.T, client *http.Client, method string, contentType, urlPath string, body []byte) (status int, b []byte, headers http.Header) {
	t.Helper()

	request, err := http.NewRequest(method, urlPath, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	return sendRequest(t, client, request, contentType)
}

func RequestWithJWT(t *testing.T, client *http.Client, jwtSecret, body []byte, method, contentType, url string, userID int) (status int, b []byte, headers http.Header) {
	t.Helper()

	jwt, err := middleware.IssueNewJWT(jwtSecret, userID)
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))

	return sendRequest(t, client, request, contentType)
}

func sendRequest(t *testing.T, client *http.Client, request *http.Request, contentType string) (status int, b []byte, headers http.Header) {
	t.Helper()

	request.Header.Set("content-type", contentType)

	resp, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	status = resp.StatusCode
	headers = resp.Header
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	b = bytes.TrimSpace(b)

	return status, b, headers
}
