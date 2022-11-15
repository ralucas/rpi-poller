package value

import (
	"time"

	"github.com/ralucas/rpi-poller/pkg/rpi"
)

type Value struct {
	Status    rpi.RPiStockStatus
	UpdatedAt time.Time
}
