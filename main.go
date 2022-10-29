package main

import (
	"log"
	"time"

	"github.com/ralucas/rpi-poller/crawler"
	"github.com/ralucas/rpi-poller/rpi"
)

func main() {
	logger := log.Default()

	logger.Println("Running rpi poller...")

	sites, err := rpi.GetSites()
	if err != nil {
		logger.Fatal(err)
	}

	c := crawler.New()

	for {
		if err := c.Crawl(sites); err != nil {
			logger.Printf("%+v", err)
		}

		time.Sleep(10 * time.Second)
	}
}
