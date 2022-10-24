package repository

import (
	"fmt"
	"sync"

	"github.com/ralucas/rpi-poller/rpi"
)

type InMemoryStore struct {
	store map[string]rpi.RPiStockStatus
	mu    sync.Mutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[string]rpi.RPiStockStatus),
	}
}

func (s *InMemoryStore) GetStockStatus(site string, productName string) (rpi.RPiStockStatus, error) {
	id := buildID(site, productName)

	s.mu.Lock()
	defer s.mu.Unlock()

	if status, found := s.store[id]; found {
		return status, nil
	}

	return -1, fmt.Errorf("item not found [%s %s]", site, productName)
}

func (s *InMemoryStore) SetStockStatus(site string, productName string, status rpi.RPiStockStatus) error {
	id := buildID(site, productName)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[id] = status

	return nil
}
