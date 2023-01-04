//go:build unit

package sms_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/ralucas/rpi-poller/internal/logging"
	"github.com/ralucas/rpi-poller/pkg/messaging/message"
	"github.com/ralucas/rpi-poller/pkg/messaging/providers/sms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	args := m.Called(url, contentType, body)

	return args.Get(0).(*http.Response), args.Error(1)
}

func TestSend(t *testing.T) {
	testurl := "test-url"
	conf := sms.Config{Url: testurl}

	testrec := "test-rec"
	testmsg := message.New("test-subj", "test-msg")

	m := struct {
		Number  string
		Message string
	}{
		Number:  testrec,
		Message: fmt.Sprintf("%s: %s", testmsg.GetSubject(), testmsg.GetMessage()),
	}

	testBody, err := json.Marshal(m)
	require.NoError(t, err)

	mockClient := &MockHttpClient{}

	mockClient.On(
		"Post", testurl, sms.ContentTypeApplicationJson, bytes.NewReader(testBody),
	).Return(
		&http.Response{StatusCode: http.StatusOK}, nil,
	)

	s := sms.New(conf, logging.NewLogger(logging.LoggerConfig{}), sms.WithHTTPClient(mockClient))

	err = s.Send(testrec, testmsg)

	mockClient.AssertNumberOfCalls(t, "Post", 1)

	assert.Nil(t, err)
}
