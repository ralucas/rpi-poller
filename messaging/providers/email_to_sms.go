package providers

import (
	"bytes"
	"fmt"
	"net/smtp"
	"net/url"

	"github.com/ralucas/rpi-poller/messaging/message"
)

type EmailToSMS struct {
	hostname string
	port     string
	username string
	password string
}

type EmailToSMSConfig struct {
	Server   string
	Username string
	Password string
}

func NewEmailToSMS(config EmailToSMSConfig) (*EmailToSMS, error) {
	u, err := url.Parse(config.Server)
	if err != nil {
		return nil, err
	}

	return &EmailToSMS{
		hostname: u.Hostname(),
		port:     u.Port(),
		username: config.Username,
		password: config.Password,
	}, nil
}

func (e *EmailToSMS) Send(msg message.Message) error {
	auth := smtp.PlainAuth("", e.username, e.password, e.hostname)
	address := fmt.Sprintf("%s:%s", e.hostname, e.port)
	from := fmt.Sprintf("%s@%s", e.username, e.hostname)

	return smtp.SendMail(address, auth, from, []string{msg.GetReceipient()}, e.messageBytes(msg))
}

func (e *EmailToSMS) messageBytes(msg message.Message) []byte {
	var bb bytes.Buffer

	bb.WriteString("To: ")
	bb.WriteString(msg.GetReceipient())
	bb.WriteString("\r\n")
	bb.WriteString("Subject: ")
	bb.WriteString(msg.GetSubject())
	bb.WriteString("\r\n\r\n")
	bb.WriteString(msg.GetMessage())
	bb.WriteString("\r\n")

	return bb.Bytes()
}
