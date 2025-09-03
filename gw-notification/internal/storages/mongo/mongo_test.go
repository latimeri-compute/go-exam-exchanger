package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/storages"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	// TODO test tables
	tests := []struct {
		name        string
		transaction storages.Transaction
		want        any
	}{
		{},
	}

	client, err := NewConnection("")
	if err != nil {
		t.Fatal(err)
	}
	m := NewWalletClient(client)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			err := m.Insert(test.transaction, ctx)
			assert.NoError(t, err)
		})
	}
}

func TestGet(t *testing.T) {
	// TODO
	// tests := []struct {
	// 	name        string
	// 	transaction storages.Transaction
	// 	want        any
	// }{
	// 	{},
	// }
}
