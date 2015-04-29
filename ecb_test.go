package ecb

import (
	"errors"
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
