package messaging

import (
	"log"

	"github.com/ralucas/rpi-poller/pkg/messaging/message"
)

type MessengerManager struct {
	listeners map[string]int
	messenger Messenger
	logger    *log.Logger
}

func NewMessengerManager(recipients []string, messenger Messenger, logging *log.Logger) *MessengerManager {
	m := &MessengerManager{
		messenger: messenger,
		listeners: make(map[string]int),
	}

	for _, r := range recipients {
		m.Subscribe(r)
	}

	return m
}

func (m *MessengerManager) Subscribe(listener string) {
	if _, ok := m.listeners[listener]; !ok {
		m.listeners[listener] = 0
	}
}

func (m *MessengerManager) Unsubsribe(listener string) {
	delete(m.listeners, listener)
}

func (m *MessengerManager) Notify(msg message.Message) error {
	errorc := make(chan error)
	for listener := range m.listeners {
		go func(l string, ch chan error) {
			ch <- m.messenger.Send(l, msg)
		}(listener, errorc)
	}

	// todo: how to handle errors?
	for err := range errorc {
		m.logger.Printf("failed sending %v", err)
	}

	return nil
}
