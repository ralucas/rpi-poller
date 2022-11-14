package providers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/ralucas/rpi-poller/messaging/message"
)

type EmailSender string

const (
	Unknown     EmailSender = "unknown"
	GmailOAuth2 EmailSender = "oauth2"
	SMTP        EmailSender = "smtp"
)

type EmailToSMS struct {
	service             *gmail.UsersMessagesService
	logger              *log.Logger
	emailSender         EmailSender
	smtpSendMail        SMTPSendMail
	credentialsFilePath string
	hostname            string
	port                string
	username            string
	password            string
}

type EmailToSMSConfig struct {
	Sender              EmailSender
	CredentialsFilePath string
	Hostname            string
	Port                string
	Username            string
	Password            string
}

type EmailToSMSOption func(*EmailToSMS)

type SMTPSendMail func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

func WithSMTPSendMailFunc(s SMTPSendMail) EmailToSMSOption {
	return func(e *EmailToSMS) {
		e.smtpSendMail = s
	}
}

func NewEmailToSMS(config EmailToSMSConfig, logger *log.Logger, options ...EmailToSMSOption) (*EmailToSMS, error) {
	var e *EmailToSMS

	switch config.Sender {
	case Unknown:
		return nil, fmt.Errorf("must build with oauth2 or smtp")
	case GmailOAuth2:
		e = &EmailToSMS{
			emailSender:         GmailOAuth2,
			credentialsFilePath: config.CredentialsFilePath,
			logger:              logger,
		}
		err := e.createService()
		if err != nil {
			return nil, fmt.Errorf("failed to create oauth2 service %v", err)
		}
	case SMTP:
		e = &EmailToSMS{
			emailSender:  SMTP,
			hostname:     config.Hostname,
			port:         config.Port,
			username:     config.Username,
			password:     config.Password,
			smtpSendMail: smtp.SendMail,
			logger:       logger,
		}
	default:
		return nil, fmt.Errorf("bad sender: %s", config.Sender)
	}

	for _, opt := range options {
		opt(e)
	}

	return e, nil
}

func (e *EmailToSMS) Send(recipient string, msg message.Message) error {
	switch e.emailSender {
	case GmailOAuth2:
		return e.sendWithGmailOAuth2(recipient, msg)
	case SMTP:
		return e.sendWithSMTP(recipient, msg)
	default:
		return fmt.Errorf("unknown sender: %s", e.emailSender)
	}
}

// GMAIL

// createService creates the gmail user messages service used for sending emails.
func (e *EmailToSMS) createService() error {
	ctx := context.Background()
	b, err := os.ReadFile(e.credentialsFilePath)
	if err != nil {
		return fmt.Errorf("failed to read client secret file: %v", err)
	}

	creds, err := google.CredentialsFromJSON(ctx, b, gmail.GmailSendScope)
	if err != nil {
		return fmt.Errorf("failed to parse credentials: %v", err)
	}

	gmailService, err := gmail.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return fmt.Errorf("failed to retrieve Gmail client: %v", err)
	}

	e.service = gmail.NewUsersMessagesService(gmailService)

	return nil
}

func (e *EmailToSMS) sendWithGmailOAuth2(recipient string, msg message.Message) error {
	sendCall := e.service.Send("me", newGmailMessage(recipient, msg))

	sentMsg, err := sendCall.Do()
	if err != nil {
		return fmt.Errorf("failed in sending the message: %v", err)
	}

	e.logger.Printf("sent message [%d]", sentMsg.ServerResponse.HTTPStatusCode)

	return nil
}

func newGmailMessage(reciepient string, msg message.Message) *gmail.Message {
	headers := []*gmail.MessagePartHeader{
		{
			Name:  "To",
			Value: reciepient,
		}, {
			Name:  "Subject",
			Value: msg.GetSubject(),
		},
	}

	body := &gmail.MessagePartBody{
		Data: msg.GetMessage(),
		Size: int64(len(msg.GetMessage())),
	}

	part := &gmail.MessagePart{
		Body:    body,
		Headers: headers,
	}

	return &gmail.Message{
		Payload: part,
	}
}

// SMTP

func (e *EmailToSMS) sendWithSMTP(recipient string, msg message.Message) error {
	auth := smtp.PlainAuth("", e.username, e.password, e.hostname)
	address := fmt.Sprintf("%s:%s", e.hostname, e.port)
	from := fmt.Sprintf("%s@%s", e.username, e.hostname)

	return e.smtpSendMail(address, auth, from, []string{recipient}, e.messageBytes(recipient, msg))
}

func (e *EmailToSMS) messageBytes(recipient string, msg message.Message) []byte {
	var bb bytes.Buffer

	bb.WriteString("To: ")
	bb.WriteString(recipient)
	bb.WriteString("\r\n")
	bb.WriteString("Subject: ")
	bb.WriteString(msg.GetSubject())
	bb.WriteString("\r\n\r\n")
	bb.WriteString(msg.GetMessage())
	bb.WriteString("\r\n")

	return bb.Bytes()
}
