package crawler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ralucas/rpi-poller/pkg/messaging/message"
	"github.com/ralucas/rpi-poller/pkg/rpi"
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

type Notifier interface {
	Notify(message.Message) error
}

type Repository interface {
	GetStockStatus(site string, productName string) (rpi.RPiStockStatus, error)
	SetStockStatus(site string, productName string, status rpi.RPiStockStatus)
}

type Crawler struct {
	logger   *log.Logger
	store    Repository
	notifier Notifier
	config   Config
}

func New(
	notifier Notifier,
	repo Repository,
	config Config,
	logger *log.Logger,
) *Crawler {
	return &Crawler{
		logger:   logger,
		store:    repo,
		notifier: notifier,
		config:   config,
	}
}

func (c *Crawler) Crawl(sites []rpi.RPiSite) error {
	errorc := make(chan error)
	resultc := make(chan []*Result)

	for _, site := range sites {
		go c.crawlSite(site, errorc, resultc)
	}

	var errors []error

	for range sites {
		select {
		case err := <-errorc:
			if err != nil {
				errors = append(errors, err)
			}
		case results := <-resultc:
			for _, result := range results {
				stockStatus := rpi.StringToStatus(result.Text)
				c.store.SetStockStatus(result.Site.Name, result.ProductName, stockStatus)

				if stockStatus == rpi.InStock {
					subject := "RPi In Stock Alert"
					msg := fmt.Sprintf("***** IN STOCK ALERT: %s - %s *****", result.Site.Name, result.ProductName)

					c.logger.Printf("sending message: %s", msg)

					err := c.notifier.Notify(message.New(subject, msg))
					if err != nil {
						c.logger.Printf("failed to send message: %+v", err)
					}
				}

				c.logger.Printf("%s - %s : %s", result.Site.Name, result.ProductName, rpi.StatusToString(stockStatus))
			}
		}
	}

	if len(errors) > 0 {
		// for now, report the first error
		return fmt.Errorf("failed to crawl %d sites [%v]", len(errors), errors[0])
	}

	return nil
}

func (c *Crawler) crawlSite(site rpi.RPiSite, errorc chan error, resultc chan []*Result) {
	// create context, ignore the initial cancel as one is given right next
	ctx, cancel1 := context.WithTimeout(context.Background(), time.Duration(c.config.TimeoutSec)*time.Second)
	defer cancel1()
	ctx, cancel2 := chromedp.NewContext(ctx)
	defer cancel2()

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
