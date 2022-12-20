package repository

import "fmt"

func BuildID(site string, productName string) string {
	return fmt.Sprintf("%s_%s", site, productName)
}
