package messaging

import (
	"fmt"

	"github.com/ralucas/rpi-poller/internal/logging"
	"github.com/ralucas/rpi-poller/pkg/messaging/message"
	"github.com/ralucas/rpi-poller/pkg/messaging/providers/emailtosms"
	"github.com/ralucas/rpi-poller/pkg/messaging/providers/sms"
)

type Messenger interface {
	Send(recipient string, msg message.Message) error
}

type Provider string

const (
	SMS        Provider = "sms"
	EmailToSMS Provider = "emailToSMS"
)

type Config struct {
	EmailToSMS emailtosms.Config
	SMS        sms.Config
}

func NewMessenger(provider Provider, config Config, logger logging.Logger) (Messenger, error) {
	switch provider {
	case SMS:
		return sms.New(config.SMS, logger), nil
	case EmailToSMS:
		return emailtosms.New(config.EmailToSMS, logger)
	}

	return nil, fmt.Errorf("no such provider: %s", provider)
}
