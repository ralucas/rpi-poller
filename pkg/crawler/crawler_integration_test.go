//go:build integration

package crawler_test

import (
	"testing"

	"github.com/ralucas/rpi-poller/internal/logging"
	"github.com/ralucas/rpi-poller/pkg/crawler"
	"github.com/ralucas/rpi-poller/pkg/crawler/mocks"
	"github.com/ralucas/rpi-poller/pkg/repository/providers/inmemory"
	"github.com/ralucas/rpi-poller/pkg/rpi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var integrationTestConfig = crawler.Config{
	BrowserTimeoutSec: 60,
	Debug:             false,
}

func isValidStatus(s rpi.RPiStockStatus) bool {
	return s == rpi.InStock || s == rpi.OutOfStock
}

func TestCrawl_Integration(t *testing.T) {
	mmm := new(mocks.MockMessengerManager)
	mmm.On("Notify", mock.Anything).Return(nil)

	l := logging.NewLogger(logging.LoggerConfig{})

	store := inmemory.New(l)

	c := crawler.New(
		mmm,
		store,
		integrationTestConfig,
		l,
	)

	testSites, err := rpi.GetSites()
	require.NoError(t, err)

	err = c.Crawl(testSites)

	assert.NoError(t, err)

	for _, site := range testSites {
		for _, prod := range site.Products {
			status, err := store.GetStockStatus(site.Name, prod.Name)
			require.NoError(t, err)

			assert.True(t, isValidStatus(status))
		}
	}
}
