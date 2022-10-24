package crawler

import (
	"context"
	"log"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/ralucas/rpi-poller/repository"
	"github.com/ralucas/rpi-poller/rpi"
)

type Result struct {
	ProductName string
	Text        string
	Ok          bool
	Attributes  map[string]string
}

type Crawler struct {
	results []*Result
	ctx     context.Context
	cancel  context.CancelFunc
	logger  *log.Logger
	store   repository.Repository
	wg      sync.WaitGroup
}

func New() *Crawler {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())

	return &Crawler{
		results: make([]*Result, 0),
		ctx:     ctx,
		cancel:  cancel,
		logger:  log.Default(),
		store:   repository.NewInMemoryStore(),
	}
}

func (c *Crawler) Start() {
	// starting browser
	if err := chromedp.Run(c.ctx); err != nil {
		c.logger.Fatalf("failed starting browser %v\n", err)
	}
}

func (c *Crawler) Crawl(sites []rpi.RPiSite) {
	defer c.cancel()

	for _, site := range sites {
		c.wg.Add(1)
		go c.crawlSite(site)
	}

	c.wg.Wait()
}

func (c *Crawler) crawlSite(site rpi.RPiSite) {
	defer c.wg.Done()

	actions := []chromedp.Action{chromedp.Navigate(site.CategoryUrl)}

	actions = append(actions, c.selectors(site.Products)...)

	c.logger.Printf("navigating to %s\n", site.CategoryUrl)

	if err := chromedp.Run(c.ctx, actions...); err != nil {
		c.logger.Fatal(err)
	}

	for _, result := range c.results {
		stockStatus := rpi.StringToStatus(result.Text)
		c.store.SetStockStatus(site.Name, result.ProductName, stockStatus)

		if stockStatus == rpi.InStock {
			// notify
			c.logger.Printf("In stock ALERT: %s - %s", site.Name, result.ProductName)
		}

		c.logger.Printf("%s - %s : %s", site.Name, result.ProductName, rpi.StatusToString(stockStatus))
	}
}

func (c *Crawler) selectors(products []rpi.RPiProduct) []chromedp.Action {
	var actions []chromedp.Action

	for _, product := range products {
		if product.Category.Attribute != "" {
			result := &Result{ProductName: product.Name}
			c.results = append(c.results, result)
			action := chromedp.AttributeValue(
				product.Category.Selector,
				product.Category.Attribute,
				&result.Text,
				&result.Ok,
				chromedp.NodeVisible,
			)
			actions = append(actions, action)
		}
	}

	return actions
}
