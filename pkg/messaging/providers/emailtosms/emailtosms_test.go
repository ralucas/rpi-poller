//go:build unit

package emailtosms_test

import (
	"fmt"
	"net/smtp"
	"testing"

	"github.com/ralucas/rpi-poller/internal/logging"

	"github.com/ralucas/rpi-poller/pkg/messaging/message"
	"github.com/ralucas/rpi-poller/pkg/messaging/providers/emailtosms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var smtpConfig = emailtosms.Config{
	Hostname: "test-hostname",
	Port:     "test-port",
	Username: "test-user",
	Password: "test-pass",
	Sender:   emailtosms.SMTP,
}

func TestEmailToSMS_New(t *testing.T) {
	e, err := emailtosms.New(smtpConfig, logging.NewLogger(logging.LoggerConfig{}))
	require.NoError(t, err)

	assert.IsType(t, &emailtosms.EmailToSMS{}, e)
}

func TestEmailToSMS_Send(t *testing.T) {
	testRec := "test-rec"
	testAddr := fmt.Sprintf("%s:%s", smtpConfig.Hostname, smtpConfig.Port)
	opt := emailtosms.WithSMTPSendMailFunc(func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		assert.Equal(t, testAddr, addr)
		assert.Equal(t, testRec, to[0])
		return nil
	})

	e, err := emailtosms.New(smtpConfig, logging.NewLogger(logging.LoggerConfig{}), opt)
	require.NoError(t, err)

	err = e.Send(testRec, message.New("test-subject", "test-msg"))
	assert.NoError(t, err)
}
