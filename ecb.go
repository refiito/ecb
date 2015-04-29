package ecb

import (
	"time"
)

type (
	Rate struct {
		Currency string
		Rate     float64
		RawRate  string
	}
	ReferenceRate struct {
		Date  time.Time
		Rates []Rate
	}
)

func RatesAt(date time.Time) (*ReferenceRate, error) {
	return rateCache.ratesAt(date)
}

func RateForAt(date time.Time, currency string) (*float64, error) {
	rate, err := RatesAt(date)
	if err != nil {
		return nil, err
	}
  if rate == nil {
    return nil, nil
  }
	return rate.RateFor(currency), nil
}

func (rate *ReferenceRate) RateFor(currency string) *float64 {
	for _, r := range rate.Rates {
		if r.Currency == currency {
			return &r.Rate
		}
	}
	return nil
}

func isSameDay(t1, t2 time.Time) (result bool) {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
