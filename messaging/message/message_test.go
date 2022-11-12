package message_test

import (
	"testing"

	"github.com/ralucas/rpi-poller/messaging/message"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	testSubj := "test-subject"
	testMsg := "test-message"
	testRec := "test-receipient"

	m := message.New(testSubj, testMsg, testRec)

	assert.Equal(t, testSubj, m.GetSubject())
	assert.Equal(t, testMsg, m.GetMessage())
	assert.Equal(t, testRec, m.GetReceipient())
}
