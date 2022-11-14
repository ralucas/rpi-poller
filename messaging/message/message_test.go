package message_test

import (
	"testing"

	"github.com/ralucas/rpi-poller/messaging/message"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	testSubj := "test-subject"
	testMsg := "test-message"

	m := message.New(testSubj, testMsg)

	assert.Equal(t, testSubj, m.GetSubject())
	assert.Equal(t, testMsg, m.GetMessage())
}
