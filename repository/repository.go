package repository

import (
	"log"

	"github.com/ralucas/rpi-poller/repository/providers"
	"github.com/ralucas/rpi-poller/rpi"
)

type Repository interface {
	GetStockStatus(site string, productName string) (rpi.RPiStockStatus, error)
	SetStockStatus(site string, productName string, status rpi.RPiStockStatus) error
}

type Provider string

const (
	InMemory Provider = "in-memory"
)

func New(provider Provider, logger *log.Logger) Repository {
	switch provider {
	case InMemory:
		return providers.NewInMemoryStore(logger)
	}

	return nil
}
