package main

import (
	"time"

	"github.com/ralucas/rpi-poller/internal/logging"
	"github.com/ralucas/rpi-poller/pkg/crawler"
	"github.com/ralucas/rpi-poller/pkg/messaging"
	"github.com/ralucas/rpi-poller/pkg/repository/providers/inmemory"
	"github.com/ralucas/rpi-poller/pkg/rpi"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger := logging.NewLogger(logging.LoggerConfig{Output: logging.FileOutput, Level: zapcore.InfoLevel})

	logger.Info("Running rpi poller...")

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

	repo := inmemory.New(logger)

	mm := messaging.NewMessengerManager(recs, m, repo, logger)

	c := crawler.New(mm, repo, conf.Crawler, logger)

	sites, err := rpi.GetSites()
	if err != nil {
		logger.Fatalf("failed to get sites %+v", err)
	}

	for {
		if err := c.Crawl(sites); err != nil {
			logger.Errorf("error crawling %+v", err)
		}

		time.Sleep(time.Duration(conf.PollTimeoutSec) * time.Second)
	}
}
