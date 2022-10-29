package messaging

import (
	"github.com/ralucas/rpi-poller/messaging/message"
	"github.com/ralucas/rpi-poller/messaging/providers"
)

type Messenger interface {
	Send(msg message.Message) error
}

type Provider string

const (
	SMS Provider = "sms"
)

func New(provider Provider) Messenger {
	switch provider {
	case SMS:
		return &providers.SMS{}
	}

	return nil
}
