package main

import (
	"log"
	"time"

	"github.com/ralucas/rpi-poller/crawler"
	"github.com/ralucas/rpi-poller/messaging"
	"github.com/ralucas/rpi-poller/rpi"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Messaging      messaging.Config
	Crawler        crawler.Config
	PollTimeoutSec int
}

func getConfig(logger *log.Logger) *AppConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatalf("failed to process config: %+v", err)
	}

	a := &AppConfig{}
	err = viper.Unmarshal(a)
	if err != nil {
		logger.Fatalf("failed to unmarshal config: %+v", err)
	}

	return a
}

func main() {
	logger := log.Default()

	logger.Println("Running rpi poller...")

	sites, err := rpi.GetSites()
	if err != nil {
		logger.Fatalf("failed to get sites: %+v", err)
	}

	config := getConfig(logger)

	c, err := crawler.New(config.Crawler, logger)
	if err != nil {
		logger.Fatalf("crawler creation failed: %+v", err)
	}

	for {
		if err := c.Crawl(sites); err != nil {
			logger.Printf("error crawling: %+v", err)
		}

		time.Sleep(time.Duration(config.PollTimeoutSec) * time.Second)
	}
}
