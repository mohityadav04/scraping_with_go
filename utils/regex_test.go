package utils

import (
	"testing"
)

func TestURLValidity(t *testing.T){
	url := "https://www.amazon.com/dp/B081ZV69YR/ref=mock"

	result := ValidateURL(url)

	if result == "" {
		t.Errorf("Not found")
	}
	
}

func TestIfPriceExist(t *testing.T){
	input_1 := "Out of Stock"
	input_2 := "$ 1.50"
	input_3 := "-1"

	result_1 := MatchPrice(input_1)
	result_2 := MatchPrice(input_2)
	result_3 := MatchPrice(input_3)

	if result_1 != "" {
		t.Errorf("Invalid Regex. Should not match if no price figure")
	}

	if result_2 == "" {
		t.Errorf("Invalid Regex. Should match if price figure exists")
	}

	if result_3 == "" {
		t.Errorf("Invalid Regex. Should match missing price label")
	}
}