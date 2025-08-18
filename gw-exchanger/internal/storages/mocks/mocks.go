package mocks

import "github.com/latimeri-compute/go-exam-exchanger/gw-exchanger/internal/storages"

type MockExchange struct {
}

func (m *MockExchange) GetAll() ([]storages.Exchange, error) {

	return nil, nil
}

func (m *MockExchange) GetRateBetween() {

}
