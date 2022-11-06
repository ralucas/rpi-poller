package messaging

import (
	"fmt"

	"github.com/ralucas/rpi-poller/messaging/message"
	"github.com/ralucas/rpi-poller/messaging/providers"
)

type Messenger interface {
	Send(msg message.Message) error
}

type Provider string

const (
	SMS Provider = "sms"
	EmailToSMS Provider = "emailToSMS"
)

type Config struct {
	providers.EmailToSMSConfig
}

func New(provider Provider, config Config) (Messenger, error) {
	switch provider {
	case SMS:
		return &providers.SMS{}, nil
	case EmailToSMS:
		return providers.NewEmailToSMS(config.EmailToSMSConfig)
	}

	return nil, fmt.Errorf("no such provider: %s", provider)
}
