package main

import (
	"log"

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

	c.Start()
	logger.Println("browser started")
	c.Crawl(sites)
}
