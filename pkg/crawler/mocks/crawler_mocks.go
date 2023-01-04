package mocks

import (
	"github.com/ralucas/rpi-poller/pkg/messaging/message"
	"github.com/ralucas/rpi-poller/pkg/rpi"
	"github.com/stretchr/testify/mock"
)

// Mock Messenger Manager

type MockMessengerManager struct {
	mock.Mock
}

func (m *MockMessengerManager) Notify(msg message.Message) error {
	args := m.Called(msg)

	return args.Error(0)
}

// Mock Repository

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetStockStatus(site string, productName string) (rpi.RPiStockStatus, error) {
	args := m.Called(site, productName)

	return args.Get(0).(rpi.RPiStockStatus), args.Error(1)
}

func (m *MockRepository) SetStockStatus(site string, productName string, status rpi.RPiStockStatus) error {
	args := m.Called(site, productName, status)

	return args.Error(0)
}
