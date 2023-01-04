package providers

import (
	"fmt"
	"sync"
	"time"

	"github.com/ralucas/rpi-poller/internal/logging"
	"github.com/ralucas/rpi-poller/pkg/model"
	"github.com/ralucas/rpi-poller/pkg/repository"
	"github.com/ralucas/rpi-poller/pkg/rpi"
)

type InMemoryStore struct {
	logger logging.Logger
	store  map[string]model.StockStatus
	notifications map[string]model.Notification
	mu     sync.Mutex
}

func NewInMemoryStore(l logging.Logger) *InMemoryStore {
	return &InMemoryStore{
		logger: l,
		store:  make(map[string]model.StockStatus),
		notifications: make(map[string]model.Notification),
	}
}

func (s *InMemoryStore) List() map[string]model.StockStatus {
	return s.store
}

func (s *InMemoryStore) GetStockStatus(site string, productName string) (rpi.RPiStockStatus, error) {
	id := repository.BuildID(site, productName)

	s.mu.Lock()
	defer s.mu.Unlock()

	if val, found := s.store[id]; found {
		return val.Status, nil
	}

	return -1, fmt.Errorf("item not found [%s %s]", site, productName)
}

func (s *InMemoryStore) SetStockStatus(site string, productName string, status rpi.RPiStockStatus) {
	id := repository.BuildID(site, productName)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[id] = model.StockStatus{Status: status, UpdatedAt: time.Now()}
}

func (s *InMemoryStore) SetNotification(recipient string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.notifications[recipient] = model.Notification{UpdatedAt: time.Now()}
}

func (s *InMemoryStore) GetNotificationByRecipient(recipient string) (notification model.Notification, exists bool) {
	val, exists := s.notifications[recipient]
	return val, exists
}