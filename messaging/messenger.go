package messaging

import (
	"fmt"
	"log"

	"github.com/ralucas/rpi-poller/messaging/message"
	"github.com/ralucas/rpi-poller/messaging/providers"
)

type Messenger interface {
	Send(msg message.Message) error
}

type Provider string

const (
	SMS        Provider = "sms"
	EmailToSMS Provider = "emailToSMS"
)

type Config struct {
	providers.EmailToSMSConfig
}

func New(provider Provider, config Config, logger *log.Logger) (Messenger, error) {
	switch provider {
	case SMS:
		return &providers.SMS{}, nil
	case EmailToSMS:
		return providers.NewEmailToSMS(config.EmailToSMSConfig, logger)
	}

	return nil, fmt.Errorf("no such provider: %s", provider)
}
