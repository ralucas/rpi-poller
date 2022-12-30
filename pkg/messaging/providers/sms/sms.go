package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ralucas/rpi-poller/pkg/messaging/message"
)

const (
	ContentTypeApplicationJson = "application/json"
)

type SMS struct {
	logger *log.Logger
	config Config
	client HttpClient
}

type Config struct {
	Url string
}

type SMSOption func(sms *SMS)

type HttpClient interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

func WithHTTPClient(client HttpClient) SMSOption {
	return func(sms *SMS) {
		sms.client = client
	}
}

func New(config Config, logger *log.Logger, opts ...SMSOption) *SMS {
	sms := &SMS{logger: logger, config: config, client: http.DefaultClient}

	for _, opt := range opts {
		opt(sms)
	}

	return sms
}

func (s *SMS) Send(recipient string, msg message.Message) error {
	m := struct {
		Number  string
		Message string
	}{
		Number:  recipient,
		Message: fmt.Sprintf("%s: %s", msg.GetSubject(), msg.GetMessage()),
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	res, err := s.client.Post(s.config.Url, ContentTypeApplicationJson, bytes.NewReader(b))
	if err != nil {
		return err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		defer res.Body.Close()

		body := make([]byte, 0)

		_, err := res.Body.Read(body)
		if err != nil {
			return err
		}

		return fmt.Errorf("failed [%d] %s", res.StatusCode, string(body))
	}

	return nil
}
