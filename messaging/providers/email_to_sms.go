package providers

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/ralucas/rpi-poller/messaging/message"
)

type EmailToSMS struct {
	service *gmail.UsersMessagesService
	logger  *log.Logger
}

type EmailToSMSConfig struct {
	CredentialsFilePath string
}

func NewEmailToSMS(config EmailToSMSConfig, logger *log.Logger) (*EmailToSMS, error) {
	svc, err := createService(config.CredentialsFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create email to sms provider: %v", err)
	}

	return &EmailToSMS{
		service: svc,
		logger:  logger,
	}, nil
}

// createService creates the gmail user messages service used for sending emails.
func createService(credentialsPath string) (*gmail.UsersMessagesService, error) {
	ctx := context.Background()
	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read client secret file: %v", err)
	}

	creds, err := google.CredentialsFromJSON(ctx, b, gmail.GmailSendScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %v", err)
	}

	client := oauth2.NewClient(ctx, creds.TokenSource)

	gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Gmail client: %v", err)
	}

	return gmail.NewUsersMessagesService(gmailService), nil
}

func (e *EmailToSMS) Send(msg message.Message) error {
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
