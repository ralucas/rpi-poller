package main

import (
	"fmt"

	"github.com/ralucas/rpi-poller/crawler"
	"github.com/ralucas/rpi-poller/messaging"
	"github.com/spf13/viper"
)

type Config struct {
	Crawler        crawler.Config
	Messaging      messaging.Config
	PollTimeoutSec int
}

func config() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to process config: %+v", err)
	}

	a := &Config{}
	err = viper.Unmarshal(a)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %+v", err)
	}

	return a, nil
}