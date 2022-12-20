package messaging

import (
	"log"
	"time"

	"github.com/ralucas/rpi-poller/pkg/messaging/message"
	"github.com/ralucas/rpi-poller/pkg/model"
)

type Store interface {
	SetNotification(recipient string)
	GetNotificationByRecipient(recipient string) (model.Notification, bool)
}

type MessengerManager struct {
	listeners map[string]int
	messenger Messenger
	store     Store
	logger    *log.Logger
}

func NewMessengerManager(recipients []string, messenger Messenger, store Store, logging *log.Logger) *MessengerManager {
	m := &MessengerManager{
		messenger: messenger,
		store:     store,
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
			if val, ok := m.store.GetNotificationByRecipient(l); ok {
				if time.Since(val.UpdatedAt).Seconds() > time.Duration.Seconds(300) {
					m.store.SetNotification(l)
					ch <- m.messenger.Send(l, msg)
				}
			} else {
				m.store.SetNotification(l)
				ch <- m.messenger.Send(l, msg)
			}
		}(listener, errorc)
	}

	// todo: how to handle errors?
	for err := range errorc {
		m.logger.Printf("failed sending %v", err)
	}

	return nil
}
