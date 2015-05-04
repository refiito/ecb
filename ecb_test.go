package ecb

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestRatesAt(t *testing.T) {
	rate, err := RatesAt(time.Date(2015, 1, 5, 12, 12, 12, 12, time.UTC))
	if err != nil {
		t.Error(err)
	}
	if rate == nil {
		t.Error(errors.New("No rate found for 5th of January 2015"))
	}
	if rate != nil && *rate.RateFor("USD") != 1.1915 {
		t.Error(errors.New("USD rate for 5th of January 2015 has changed..."))
	}
}

func TestRateForAt(t *testing.T) {
	rate, err := RateForAt(time.Date(2015, 1, 5, 12, 12, 12, 12, time.UTC), "USD")
	if err != nil {
		t.Error(err)
	}
	if rate == nil {
		t.Error(errors.New("No rate found"))
	}
	if rate != nil && *rate != 1.1915 {
		t.Error(errors.New("USD rate for 5th of January 2015 has changed..."))
	}
}

func TestRatesForBetween(t *testing.T) {
	startTime := time.Date(2015, 4, 5, 12, 12, 12, 0, time.UTC)
	endTime := time.Date(2015, 4, 15, 12, 12, 12, 0, time.UTC)
	rates, err := RatesForBetween(startTime, endTime, "USD")
	if err != nil {
		t.Error(err)
	}
	if len(rates) != 7 {
		t.Error(errors.New("Got different number of rates than expected"))
	}
	if !rates[0].Date.Equal(time.Date(2015, 4, 15, 0, 0, 0, 0, time.UTC)) {
		t.Error(fmt.Errorf("Element 0 has different time than expected, %v", rates[0].Date))
	}
	if !rates[1].Date.Equal(time.Date(2015, 4, 14, 0, 0, 0, 0, time.UTC)) {
		t.Error(fmt.Errorf("Element 1 has different time than expected, %v", rates[0].Date))
	}
	if !rates[2].Date.Equal(time.Date(2015, 4, 13, 0, 0, 0, 0, time.UTC)) {
		t.Error(fmt.Errorf("Element 2 has different time than expected, %v", rates[0].Date))
	}
	if !rates[3].Date.Equal(time.Date(2015, 4, 10, 0, 0, 0, 0, time.UTC)) {
		t.Error(fmt.Errorf("Element 3 has different time than expected, %v", rates[0].Date))
	}
	if !rates[4].Date.Equal(time.Date(2015, 4, 9, 0, 0, 0, 0, time.UTC)) {
		t.Error(fmt.Errorf("Element 4 has different time than expected, %v", rates[0].Date))
	}
	if !rates[5].Date.Equal(time.Date(2015, 4, 8, 0, 0, 0, 0, time.UTC)) {
		t.Error(fmt.Errorf("Element 5 has different time than expected, %v", rates[0].Date))
	}
	if !rates[6].Date.Equal(time.Date(2015, 4, 7, 0, 0, 0, 0, time.UTC)) {
		t.Error(fmt.Errorf("Element 6 has different time than expected, %v", rates[0].Date))
	}
}
