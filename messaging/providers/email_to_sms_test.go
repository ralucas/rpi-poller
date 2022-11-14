package providers_test

import (
	"fmt"
	"log"
	"net/smtp"
	"testing"

	"github.com/ralucas/rpi-poller/messaging/message"
	"github.com/ralucas/rpi-poller/messaging/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var smtpConfig = providers.EmailToSMSConfig{
	Hostname: "test-hostname",
	Port:     "test-port",
	Username: "test-port",
	Password: "test-pass",
	Sender:   providers.SMTP,
}

func TestEmailToSMS_New(t *testing.T) {
	e, err := providers.NewEmailToSMS(smtpConfig, log.Default())
	require.NoError(t, err)

	assert.IsType(t, &providers.EmailToSMS{}, e)
}

func TestEmailToSMS_Send(t *testing.T) {
	testRec := "test-rec"
	testAddr := fmt.Sprintf("%s:%s", smtpConfig.Hostname, smtpConfig.Port)
	opt := providers.WithSMTPSendMailFunc(func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		assert.Equal(t, testAddr, addr)
		assert.Equal(t, testRec, to[0])
		return nil
	})

	e, err := providers.NewEmailToSMS(smtpConfig, log.Default(), opt)
	require.NoError(t, err)

	err = e.Send(testRec, message.New("test-subject", "test-msg"))
	assert.NoError(t, err)
}
