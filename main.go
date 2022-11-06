package main

import (
	"flag"
	"log"
	"time"

	"github.com/ralucas/rpi-poller/crawler"
	"github.com/ralucas/rpi-poller/rpi"
	
	"github.com/kelseyhightower/envconfig"
)

var pollTimeout int

func init() {
	flag.IntVar(&pollTimeout, "t", 10, "poll timeout (in seconds)")
}

func main() {
	logger := log.Default()

	logger.Println("Running rpi poller...")

	err := envconfig.Process("myapp", &s)

	sites, err := rpi.GetSites()
	if err != nil {
		logger.Fatal(err)
	}

	c, err := crawler.New()
	if err != nil {
		logger.Fatal(err)
	}

	for {
		if err := c.Crawl(sites); err != nil {
			logger.Printf("%+v", err)
		}

		time.Sleep(time.Duration(pollTimeout) * time.Second)
	}
}
