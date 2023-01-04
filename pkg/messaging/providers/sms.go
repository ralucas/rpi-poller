package providers

import (
	"github.com/ralucas/rpi-poller/internal/logging"
	"github.com/ralucas/rpi-poller/pkg/messaging/message"
)

type SMS struct {
	logger logging.Logger
}

func NewSMS(logger logging.Logger) *SMS {
	return &SMS{logger: logger}
}

func (s *SMS) Send(recipient string, msg message.Message) error {
	return nil
}
