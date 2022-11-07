package crawler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ralucas/rpi-poller/messaging"
	"github.com/ralucas/rpi-poller/messaging/message"
	"github.com/ralucas/rpi-poller/repository"
	"github.com/ralucas/rpi-poller/rpi"
)

const TIMEOUT_SEC = 10

type Result struct {
	ProductName string
	Text        string
	Ok          bool
	Attributes  map[string]string
}

type Crawler struct {
	logger    *log.Logger
	store     repository.Repository
	messenger messaging.Messenger
}

func New(logger *log.Logger) (*Crawler, error) {
	m, err := messaging.New(messaging.EmailToSMS, messaging.Config{}, logger)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		logger:    logger,
		store:     repository.New(repository.InMemory, logger),
		messenger: m,
	}, nil
}

func (c *Crawler) Crawl(sites []rpi.RPiSite) error {
	var cancel context.CancelFunc
	var ctx context.Context

	errorc := make(chan error)

	for _, site := range sites {
		ctx, cancel = context.WithCancel(context.Background())
		go c.crawlSite(ctx, site, errorc)
	}

	defer cancel()

	for range sites {
		if err := <-errorc; err != nil {
			cancel()
			return err
		}
	}

	return nil
}

func (c *Crawler) crawlSite(ctx context.Context, site rpi.RPiSite, errorc chan error) {
	// create context
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// starting browser
	if err := chromedp.Run(ctx); err != nil {
		c.logger.Fatalf("failed starting browser %v\n", err)
		errorc <- err
	}

	ctx, cancel = context.WithTimeout(ctx, TIMEOUT_SEC*time.Second)
	defer cancel()

	actions := []chromedp.Action{chromedp.Navigate(site.CategoryUrl)}

	selActions, results := c.selectors(site.Products)

	actions = append(actions, selActions...)

	c.logger.Printf("navigating to %s\n", site.CategoryUrl)

	if err := chromedp.Run(ctx, actions...); err != nil {
		errorc <- fmt.Errorf("failed crawling %s: %+v", site.CategoryUrl, err)
		return
	}

	for _, result := range results {
		stockStatus := rpi.StringToStatus(result.Text)
		c.store.SetStockStatus(site.Name, result.ProductName, stockStatus)

		if stockStatus == rpi.InStock {
			subject := "RPi In Stock Alert"
			msg := fmt.Sprintf("***** IN STOCK ALERT: %s - %s *****", site.Name, result.ProductName)
			c.logger.Printf(msg)
			err := c.messenger.Send(message.New(subject, msg, ""))
			if err != nil {
				c.logger.Printf("failed to send message: %+v", err)
			}
		}

		c.logger.Printf("%s - %s : %s", site.Name, result.ProductName, rpi.StatusToString(stockStatus))
	}

	errorc <- nil
}

func (c *Crawler) selectors(products []rpi.RPiProduct) ([]chromedp.Action, []*Result) {
	var actions []chromedp.Action
	var results []*Result

	for _, product := range products {
		if product.Category.Attribute != "" {
			result := &Result{ProductName: product.Name}
			results = append(results, result)
			action := chromedp.AttributeValue(
				product.Category.Selector,
				product.Category.Attribute,
				&result.Text,
				&result.Ok,
			)
			actions = append(actions, action)
		} else {
			result := &Result{ProductName: product.Name}
			results = append(results, result)
			action := chromedp.Text(
				product.Category.Selector,
				&result.Text,
			)
			actions = append(actions, action)

		}
	}

	return actions, results
}
