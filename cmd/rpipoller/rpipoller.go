package main

import (
	"log"
	"time"

	"github.com/ralucas/rpi-poller/pkg/crawler"
	"github.com/ralucas/rpi-poller/pkg/messaging"
	"github.com/ralucas/rpi-poller/pkg/repository"
	"github.com/ralucas/rpi-poller/pkg/rpi"
)

func main() {
	logger := log.Default()

	logger.Println("Running rpi poller...")

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

	repo := repository.New(repository.InMemory, logger)

	c := crawler.New(mm, repo, conf.Crawler, logger)

	sites, err := rpi.GetSites()
	if err != nil {
		logger.Fatalf("failed to get sites %+v", err)
	}

	for {
		if err := c.Crawl(sites); err != nil {
			logger.Printf("error crawling %+v", err)
		}

		time.Sleep(time.Duration(conf.PollTimeoutSec) * time.Second)
	}
}
