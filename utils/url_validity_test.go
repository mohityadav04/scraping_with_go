package utils

import (
	"testing"
)

func TestURLValidity(t *testing.T){
	url := "https://www.amazon.com/dp/B081ZV69YR/ref=mock"

	result, err := ValidateURL(url)

	if err != nil {
		t.Errorf("Failed to parse")
	}
	if result == false {
		t.Errorf("wrong result")
	}
	
}