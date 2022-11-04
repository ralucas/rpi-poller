package main

import (
	"flag"
	"log"
	"time"

	"github.com/ralucas/rpi-poller/crawler"
	"github.com/ralucas/rpi-poller/rpi"
)

var pollTimeout int

func init() {
	flag.IntVar(&pollTimeout, "t", 10, "poll timeout (in seconds)")
}

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

		time.Sleep(time.Duration(pollTimeout) * time.Second)
	}
}
