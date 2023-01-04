package crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ralucas/rpi-poller/internal/logging"
	"github.com/ralucas/rpi-poller/pkg/messaging/message"
	"github.com/ralucas/rpi-poller/pkg/rpi"
)

type Config struct {
	BrowserTimeoutSec int
	Debug             bool
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
	SetStockStatus(site string, productName string, status rpi.RPiStockStatus) error
}

type Crawler struct {
	logger   logging.Logger
	store    Repository
	notifier Notifier
	config   Config
}

func New(
	notifier Notifier,
	repo Repository,
	config Config,
	logger logging.Logger,
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
	resultc := make(chan *Result)

	n := 0

	for _, site := range sites {
		go c.crawlSite(site, errorc, resultc)
		n += len(site.Products)
	}

	var errors []error

	for i := 0; i < n; i++ {
		select {
		case err := <-errorc:
			if err != nil {
				errors = append(errors, err)
				c.logger.Errorf(err.Error())
			}
		case result := <-resultc:
			stockStatus := rpi.StringToStatus(result.Text)

			err := c.store.SetStockStatus(result.Site.Name, result.ProductName, stockStatus)
			if err != nil {
				errors = append(errors, err)
				c.logger.Errorf(err.Error())
				continue
			}

			if stockStatus == rpi.InStock {
				subject := "RPi In Stock Alert"
				msg := fmt.Sprintf("***** IN STOCK ALERT: %s - %s *****", result.Site.Name, result.ProductName)

				c.logger.Infof("sending message: %s", msg)

				err := c.notifier.Notify(message.New(subject, msg))
				if err != nil {
					c.logger.Infof("failed to send message: %+v", err)
				}
			}

			c.logger.Infof("%s - %s : %s [%s]", result.Site.Name, result.ProductName, rpi.StatusToString(stockStatus), result.Text)
		}
	}

	if len(errors) > 0 {
		// for now, report the first error
		return fmt.Errorf("failed to crawl %d products [%v]", len(errors), errors[0])
	}

	return nil
}

func (c *Crawler) crawlSite(site rpi.RPiSite, errorc chan error, resultc chan *Result) {
	var ctx context.Context

	// create context, ignore the initial cancel as one is given right next
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.config.BrowserTimeoutSec)*time.Second)
	defer cancel()

	if c.config.Debug {
		var cancel1 context.CancelFunc
		ctx, cancel1 = chromedp.NewExecAllocator(ctx, append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
		defer cancel1()
	}

	ctx, cancel2 := chromedp.NewContext(ctx)
	defer cancel2()

	for _, product := range site.Products {
		// starting tab
		if err := chromedp.Run(ctx); err != nil {
			c.logger.Fatalf("failed starting browser %v\n", err)
			errorc <- err
			return
		}
		// Results are attached to the selectors and bound as pointers
		// so they get passed back here as pointers to be populated
		// by the `Run` process later and returned via a channel.
		result := &Result{Site: site, ProductName: product.Name}
		actions, result := c.createActions(product, result)

		c.logger.Infof("navigating to %s\n", product.Url)

		if err := chromedp.Run(ctx, actions...); err != nil {
			errorc <- fmt.Errorf("failed crawling %s for %s: %+v", product.Url, product.Name, err)
			return
		}

		var cancel3 context.CancelFunc
		ctx, cancel3 = chromedp.NewContext(ctx)
		defer cancel3()

		resultc <- result
	}
}

func (c *Crawler) createActions(product rpi.RPiProduct, result *Result) ([]chromedp.Action, *Result) {
	actions := []chromedp.Action{chromedp.Navigate(product.Url)}

	if product.Attribute != "" {
		action := chromedp.AttributeValue(
			product.Selector,
			product.Attribute,
			&result.Text,
			&result.Ok,
		)
		actions = append(actions, action)
	} else {
		action := chromedp.Text(
			product.Selector,
			&result.Text,
		)
		actions = append(actions, action)
	}

	return actions, result
}
