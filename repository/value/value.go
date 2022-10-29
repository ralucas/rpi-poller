package value

import (
	"time"

	"github.com/ralucas/rpi-poller/rpi"
)

type Value struct {
	Status    rpi.RPiStockStatus
	UpdatedAt time.Time
}
