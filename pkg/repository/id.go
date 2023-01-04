package repository

import "fmt"

func BuildID(site string, productName string) (string, error) {
	if site == "" || productName == "" {
		return "", fmt.Errorf("requires site [%s] and product name [%s]", site, productName)
	}
	return fmt.Sprintf("%s_%s", site, productName), nil
}
