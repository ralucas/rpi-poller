package main

import (
	"log"
	"time"

	"github.com/ralucas/rpi-poller/crawler"
	"github.com/ralucas/rpi-poller/messaging"
	"github.com/ralucas/rpi-poller/rpi"
)

func main() {
	logger := log.Default()

	logger.Println("Running rpi poller...")

	sites, err := rpi.GetSites()
	if err != nil {
		logger.Fatalf("failed to get sites %+v", err)
	}

	conf, err := config()
	if err != nil {
		logger.Fatalf("failed to get config %v", err)
	}

	m, err := messaging.NewMessenger(messaging.EmailToSMS, conf.Messaging, logger)
	if err != nil {
		logger.Fatalf("failed to create messenger %v", err)
	}

	recs, err := recipients()
	if err != nil {
		logger.Fatalf("failed to get recipients %v", err)
	}

	mm := messaging.NewMessengerManager(recs, m, logger)

	c, err := crawler.New(mm, conf.Crawler, logger)
	if err != nil {
		logger.Fatalf("failed to create crawler %+v", err)
	}

	for {
		if err := c.Crawl(sites); err != nil {
			logger.Printf("error crawling %+v", err)
		}

		time.Sleep(time.Duration(conf.PollTimeoutSec) * time.Second)
	}
}
