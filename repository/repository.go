package repository

import (
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

func New(provider Provider) Repository {
	switch provider {
	case InMemory:
		return providers.NewInMemoryStore()
	}

	return nil
}
