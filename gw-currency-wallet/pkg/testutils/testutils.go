package testutils

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func HashString(t *testing.T, str string) []byte {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	return bytes
}

func ReceiveResponseBody(t *testing.T, client *http.Client, method string, urlPath string, body []byte) []byte {
	t.Helper()

	req, err := http.NewRequest(method, urlPath, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
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
