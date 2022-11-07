package providers

import (
	"log"

	"github.com/ralucas/rpi-poller/messaging/message"
)

type SMS struct {
	logger *log.Logger
}

func (s *SMS) Send(msg message.Message) error {
	return nil
}
