package repository

import (
	"fmt"

	"github.com/ralucas/rpi-poller/rpi"
)

type Repository interface {
	GetStockStatus(site string, productName string) (rpi.RPiStockStatus, error)
	SetStockStatus(site string, productName string, status rpi.RPiStockStatus) error
}

func buildID(site string, productName string) string {
	return fmt.Sprintf("%s_%s", site, productName)
}
