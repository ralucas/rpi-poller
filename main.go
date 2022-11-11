package main

import (
	"log"
	"time"

	"github.com/ralucas/rpi-poller/crawler"
	"github.com/ralucas/rpi-poller/messaging"
	"github.com/ralucas/rpi-poller/rpi"

	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	messaging.Config
	PollTimeout int
}

func main() {
	logger := log.Default()

	logger.Println("Running rpi poller...")

	a := AppConfig{}
	err := envconfig.Process("", &a)
	if err != nil {
		logger.Fatalf("failed to process config: %+v", err)
	}

	sites, err := rpi.GetSites()
	if err != nil {
		logger.Fatalf("failed to get sites: %+v", err)
	}

	c, err := crawler.New(logger)
	if err != nil {
		logger.Fatalf("crawler creation failed: %+v", err)
	}

	for {
		if err := c.Crawl(sites); err != nil {
			logger.Printf("%+v", err)
		}

		time.Sleep(time.Duration(a.PollTimeout) * time.Second)
	}
}
