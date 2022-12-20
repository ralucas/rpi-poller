package crawler_test

import (
	"log"
	"testing"

	"github.com/ralucas/rpi-poller/pkg/crawler"
	"github.com/ralucas/rpi-poller/pkg/crawler/mocks"
	"github.com/ralucas/rpi-poller/pkg/rpi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testConfig = crawler.Config{
	TimeoutSec: 10,
}

var testSites = []rpi.RPiSite{{
	Name:        "test-name",
	CategoryUrl: "https://yahoo.com",
	Products: []rpi.RPiProduct{{
		Name: "test-product",
		Url:  "https://yahoo.com",
		Category: rpi.RPiProductCategory{
			Selector:  "selc",
			Attribute: "attr",
		},
	}},
}}

func TestCrawl(t *testing.T) {
	mmm := new(mocks.MockMessengerManager)
	mmm.On("Notify", mock.Anything).Return(nil)

	mr := new(mocks.MockRepository)
	mr.On("SetStockStatus", mock.Anything, mock.Anything, mock.Anything)
	mr.On("GetStockStatus").Return(0, nil)

	c := crawler.New(
		mmm,
		mr,
		testConfig,
		log.Default(),
	)

	err := c.Crawl(testSites)

	assert.NoError(t, err)
}
