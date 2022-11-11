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

type Config struct {
	TimeoutSec int
}

type Result struct {
	Site        rpi.RPiSite
	ProductName string
	Text        string
	Ok          bool
	Attributes  map[string]string
}

type Crawler struct {
	logger    *log.Logger
	store     repository.Repository
	messenger messaging.Messenger
	config    Config
}

func New(config Config, logger *log.Logger) (*Crawler, error) {
	m, err := messaging.New(messaging.EmailToSMS, messaging.Config{}, logger)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		logger:    logger,
		store:     repository.New(repository.InMemory, logger),
		messenger: m,
		config:    config,
	}, nil
}

func (c *Crawler) Crawl(sites []rpi.RPiSite) error {
	errorc := make(chan error)
	resultc := make(chan []*Result)

	for _, site := range sites {
		go c.crawlSite(site, errorc, resultc)
	}

	for range sites {
		if err := <-errorc; err != nil {
			return err
		}

		for _, result := range <-resultc {
			stockStatus := rpi.StringToStatus(result.Text)
			c.store.SetStockStatus(result.Site.Name, result.ProductName, stockStatus)

			if stockStatus == rpi.InStock {
				subject := "RPi In Stock Alert"
				msg := fmt.Sprintf("***** IN STOCK ALERT: %s - %s *****", result.Site.Name, result.ProductName)

				c.logger.Printf("sending message: %s", msg)

				err := c.messenger.Send(message.New(subject, msg, ""))
				if err != nil {
					c.logger.Printf("failed to send message: %+v", err)
				}
			}

			c.logger.Printf("%s - %s : %s", result.Site.Name, result.ProductName, rpi.StatusToString(stockStatus))
		}
	}

	return nil
}

func (c *Crawler) crawlSite(site rpi.RPiSite, errorc chan error, resultc chan []*Result) {
	// create context, ignore the initial cancel as one is given right next
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(c.config.TimeoutSec)*time.Second)
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// starting browser
	if err := chromedp.Run(ctx); err != nil {
		c.logger.Fatalf("failed starting browser %v\n", err)
		errorc <- err
		return
	}

	actions := []chromedp.Action{chromedp.Navigate(site.CategoryUrl)}

	// Results are attached to the selectors and bound as pointers
	// so they get passed back here as pointers to be populated
	// by the `Run` process later and returned via a channel.
	selActions, results := c.selectors(site.Products)

	actions = append(actions, selActions...)

	c.logger.Printf("navigating to %s\n", site.CategoryUrl)

	if err := chromedp.Run(ctx, actions...); err != nil {
		errorc <- fmt.Errorf("failed crawling %s: %+v", site.CategoryUrl, err)
		return
	}

	resultc <- results
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
