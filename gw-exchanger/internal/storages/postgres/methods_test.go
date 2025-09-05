package postgres

import (
	"testing"

	"github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages"
	"github.com/stretchr/testify/assert"
)

var dsn = "postgresql://postgres:password@localhost:5432/exchanger_test"

func TestGetAll(t *testing.T) {
	want := []storages.ReturnExchanges{
		{
			FromValuteCode: "usd",
			ToValuteCode:   "rub",
			Rate:           803466,
		},
		{
			FromValuteCode: "eur",
			ToValuteCode:   "rub",
			Rate:           935604,
		},
		{
			FromValuteCode: "rub",
			ToValuteCode:   "usd",
			Rate:           124,
		},
		{
			FromValuteCode: "rub",
			ToValuteCode:   "eur",
			Rate:           107,
		},
		{
			FromValuteCode: "usd",
			ToValuteCode:   "eur",
			Rate:           8588,
		},
		{
			FromValuteCode: "eur",
			ToValuteCode:   "usd",
			Rate:           11645,
		},
	}

	db, err := NewConnection(dsn)
	if err != nil {
		t.Fatal(err)
	}
	setupDB(t, db)
	defer teardownDB(t, db)

	rates, err := db.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rates)
	assert.EqualValues(t, want, rates)
}

func TestGetRateBetween(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		wantRes any
		wantErr error
	}{
		{
			from: "rub",
			to:   "usd",
			wantRes: storages.ReturnExchanges{
				FromValuteCode: "rub",
				ToValuteCode:   "usd",
				Rate:           124,
			},
		},
		{
			from: "usd",
			to:   "rub",
			wantRes: storages.ReturnExchanges{
				FromValuteCode: "usd",
				ToValuteCode:   "rub",
				Rate:           803466,
			},
		},
		{
			from:    "awawa",
			to:      "rub",
			wantRes: storages.ReturnExchanges{},
			wantErr: storages.ErrNotFound,
		},
	}
	db, err := NewConnection(dsn)
	if err != nil {
		t.Fatal(err)
	}
	setupDB(t, db)
	defer teardownDB(t, db)

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			got, err := db.GetRateBetween(test.from, test.to)
			assert.ErrorIs(t, err, test.wantErr)
			assert.Equal(t, test.wantRes, got)
		})
	}
}
