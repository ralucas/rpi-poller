package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Recipients struct {
	Recipients []string
}

func recipients() ([]string, error) {
	recYaml, err := os.ReadFile("./data/recipients.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read yaml: %v ", err)
	}

	r := &Recipients{}
	err = yaml.Unmarshal(recYaml, r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	return r.Recipients, nil
}
