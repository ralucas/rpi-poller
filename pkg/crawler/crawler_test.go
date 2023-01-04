//go:build unit

package crawler_test

import (
	"testing"

	"github.com/ralucas/rpi-poller/internal/logging"
	"github.com/ralucas/rpi-poller/pkg/crawler"
	"github.com/ralucas/rpi-poller/pkg/crawler/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testConfig = crawler.Config{
	BrowserTimeoutSec: 60,
	Debug:             false,
}

func TestNew(t *testing.T) {
	mmm := new(mocks.MockMessengerManager)
	mmm.On("Notify", mock.Anything).Return(nil)

	mr := new(mocks.MockRepository)
	mr.On("SetStockStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mr.On("GetStockStatus").Return(0, nil)

	l := logging.NewLogger(logging.LoggerConfig{})

	c := crawler.New(
		mmm,
		mr,
		testConfig,
		l,
	)

	assert.IsType(t, &crawler.Crawler{}, c)
}
