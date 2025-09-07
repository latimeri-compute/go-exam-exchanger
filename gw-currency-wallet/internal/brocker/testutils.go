package brocker

import (
	"testing"

	"github.com/IBM/sarama/mocks"
)

func NewTestProducer(t *testing.T, mock *mocks.SyncProducer) *Producer {

	return &Producer{mock}
}
