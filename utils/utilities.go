package utils

import (
	"regexp"
)

func ValidateURL(url string) string {
	re := regexp.MustCompile("(/dp/[A-Z0-9]{10})")
	return re.FindString(url)
}

func MatchPrice(url string) string {
	re := regexp.MustCompile("[0-9]+")
	return re.FindString(url)
}