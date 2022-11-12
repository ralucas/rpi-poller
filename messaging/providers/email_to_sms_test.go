package providers_test

import (
	"log"
	"testing"

	"github.com/ralucas/rpi-poller/messaging/message"
	"github.com/ralucas/rpi-poller/messaging/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const credsPath = "../../test/resources/credentials.json"

func TestEmailToSMS_New(t *testing.T) {
	builder := providers.NewEmailToSMSBuilder(log.Default())
	builder.WithSMTP("test-hostname", "test-port", "test-username", "test-password")
	e, err := builder.Build()
	require.NoError(t, err)

	assert.IsType(t, &providers.EmailToSMS{}, e)
}

func TestEmailToSMS_Send(t *testing.T) {
	builder := providers.NewEmailToSMSBuilder(log.Default())
	builder.WithSMTP("test-hostname", "test-port", "test-username", "test-password")
	e, err := builder.Build()
	require.NoError(t, err)

	err = e.Send(message.New("test-subject", "test-msg", "test-rec"))
	assert.Error(t, err)
}
