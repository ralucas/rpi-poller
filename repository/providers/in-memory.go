package providers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ralucas/rpi-poller/repository/util"
	"github.com/ralucas/rpi-poller/repository/value"
	"github.com/ralucas/rpi-poller/rpi"
)

type InMemoryStore struct {
	logger *log.Logger
	store  map[string]value.Value
	mu     sync.Mutex
}

func NewInMemoryStore(l *log.Logger) *InMemoryStore {
	return &InMemoryStore{
		logger: l,
		store:  make(map[string]value.Value),
	}
}

func (s *InMemoryStore) List() map[string]value.Value {
	return s.store
}

func (s *InMemoryStore) GetStockStatus(site string, productName string) (rpi.RPiStockStatus, error) {
	id := util.BuildID(site, productName)

	s.mu.Lock()
	defer s.mu.Unlock()

	if val, found := s.store[id]; found {
		return val.Status, nil
	}

	return -1, fmt.Errorf("item not found [%s %s]", site, productName)
}

func (s *InMemoryStore) SetStockStatus(site string, productName string, status rpi.RPiStockStatus) error {
	id := util.BuildID(site, productName)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[id] = value.Value{Status: status, UpdatedAt: time.Now()}

	return nil
}
