package model

import (
	"time"

	"github.com/ralucas/rpi-poller/pkg/rpi"
)

type StockStatus struct {
	Status    rpi.RPiStockStatus
	UpdatedAt time.Time
}
