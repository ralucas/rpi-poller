package providers

import (
	"log"

	"github.com/ralucas/rpi-poller/pkg/messaging/message"
)

type SMS struct {
	logger *log.Logger
}

func NewSMS(logger *log.Logger) *SMS {
	return &SMS{logger: logger}
}

func (s *SMS) Send(recipient string, msg message.Message) error {
	return nil
}
