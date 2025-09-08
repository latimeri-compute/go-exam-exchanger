package brocker

import (
	"testing"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-notification/internal/storages"
	"github.com/stretchr/testify/assert"
)

// func TestMain(m *testing.M) {
// 	sarama.Logger = log.New(os.Stdout, "[testing]", log.LstdFlags)
// 	m.Run()
// }

// TODO
func TestConsumer(t *testing.T) {

	ch := make(chan storages.Transaction)
	group := newTestConsumers(t, ch)
	defer group.Group.Close()
	prod := newProducer(t)

	tests := []struct {
		name        string
		transaction transactionMessage
		want        storages.Transaction
	}{
		{
			name: "exchange",
			transaction: transactionMessage{
				WalletID:     1,
				Type:         "exchange",
				FromCurrency: "rub",
				ToCurrency:   "usd",
				AmountFrom:   300000,
				AmountTo:     200000,
			},
			want: storages.Transaction{
				WalletID:     1,
				Type:         "exchange",
				FromCurrency: "rub",
				ToCurrency:   "usd",
				AmountFrom:   300000,
				AmountTo:     200000,
			},
		},
		{
			name: "withdraw",
			transaction: transactionMessage{
				WalletID:     1,
				Type:         "withdraw",
				FromCurrency: "rub",
				AmountFrom:   300000,
			},
			want: storages.Transaction{
				WalletID:     1,
				Type:         "withdraw",
				FromCurrency: "rub",
				AmountFrom:   300000,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tm := time.Now()
			test.transaction.Timestamp = tm
			test.want.Timestamp = tm
			send(t, prod, test.transaction)
			res := <-ch

			assert.EqualExportedValues(t, test.want, res)
		})
	}
}
