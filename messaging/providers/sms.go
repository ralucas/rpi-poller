package providers

import "github.com/ralucas/rpi-poller/messaging/message"

type SMS struct {
}

func (s *SMS) Send(msg message.Message) error {
	return nil
}
