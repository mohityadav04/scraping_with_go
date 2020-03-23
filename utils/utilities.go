package utils

import (
	"regexp"
)

func ValidateURL(url string) (bool, error) {
	return regexp.MatchString(`(/dp/[A-Z0-9]{10}/)`, url)
}