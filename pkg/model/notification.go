package model

import (
	"time"

	"github.com/ralucas/rpi-poller/pkg/messaging/message"
)

type Notification struct {
	Msg       message.Message
	UpdatedAt time.Time
}
