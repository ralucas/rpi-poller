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
	sender              EmailSender
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

type EmailToSMSBuilder struct {
	e *EmailToSMS
}

func (b *EmailToSMSBuilder) WithOauth2(credentialsFilePath string) *EmailToSMSBuilder {
	b.e.sender = GmailOAuth2
	b.e.credentialsFilePath = credentialsFilePath

	return b
}

func (b *EmailToSMSBuilder) WithSMTP(hostname, port, username, password string) *EmailToSMSBuilder {
	b.e.sender = SMTP
	b.e.hostname = hostname
	b.e.port = port
	b.e.username = username
	b.e.password = password

	return b
}

func (b *EmailToSMSBuilder) Build() (*EmailToSMS, error) {
	if b.e.sender == Unknown {
		return nil, fmt.Errorf("must build with oauth2 or smtp")
	}

	if b.e.sender == GmailOAuth2 {
		err := b.e.createService()
		if err != nil {
			return nil, fmt.Errorf("failed to create oauth2 service %v", err)
		}
	}

	return b.e, nil
}

func NewEmailToSMSBuilder(logger *log.Logger) *EmailToSMSBuilder {
	e := &EmailToSMS{
		logger: logger,
	}

	return &EmailToSMSBuilder{e}
}

func (e *EmailToSMS) Send(msg message.Message) error {
	switch e.sender {
	case GmailOAuth2:
		return e.sendWithGmailOAuth2(msg)
	case SMTP:
		return e.sendWithSMTP(msg)
	default:
		return fmt.Errorf("unknown sender: %s", e.sender)
	}
}

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

func (e *EmailToSMS) sendWithGmailOAuth2(msg message.Message) error {
	sendCall := e.service.Send("me", newGmailMessage(msg))

	sentMsg, err := sendCall.Do()
	if err != nil {
		return fmt.Errorf("failed in sending the message: %v", err)
	}

	e.logger.Printf("sent message [%d]", sentMsg.ServerResponse.HTTPStatusCode)

	return nil
}

func newGmailMessage(msg message.Message) *gmail.Message {
	headers := []*gmail.MessagePartHeader{
		{
			Name:  "To",
			Value: msg.GetReceipient(),
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

func (e *EmailToSMS) sendWithSMTP(msg message.Message) error {
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
