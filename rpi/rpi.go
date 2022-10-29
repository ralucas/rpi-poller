package rpi

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type RPi struct {
	Sites []RPiSite
}

type RPiSite struct {
	Name        string `yaml:"name"`
	CategoryUrl string `yaml:"category_url"`
	Products    []RPiProduct
}

type RPiProduct struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	Category struct {
		Selector  string `yaml:"selector"`
		Attribute string `yaml:"attribute"`
	}
}

type RPiStockStatus int

const (
	OutOfStock RPiStockStatus = iota
	InStock
	Unknown
)

func GetSites() ([]RPiSite, error) {
	rpi := &RPi{}

	rpiYaml, err := os.ReadFile("rpi.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read yaml: %v ", err)
	}

	err = yaml.Unmarshal(rpiYaml, rpi)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	return rpi.Sites, nil
}

func StringToStatus(s string) RPiStockStatus {
	prepared := strings.ToLower(s)

	if strings.Contains(prepared, "out") {
		return OutOfStock
	} else if strings.Contains(prepared, "in") {
		return InStock
	}

	return Unknown
}

func StatusToString(s RPiStockStatus) string {
	switch s {
	case OutOfStock:
		return "Out of Stock"
	case InStock:
		return "In Stock"
	default:
		return "unknown status"
	}
}
