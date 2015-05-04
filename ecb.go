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
	CurrencyRate struct {
		Date     time.Time
		Currency string
		Rate     *float64
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

func RatesForBetween(rangeStart, rangeEnd time.Time, currency string) ([]CurrencyRate, error) {
	var result []CurrencyRate
	rates, err := rateCache.ratesBetween(rangeStart, rangeEnd)
	if err != nil {
		return result, err
	}
	for _, rate := range rates {
		result = append(result, CurrencyRate{Date: rate.Date, Currency: currency, Rate: rate.RateFor(currency)})
	}
	return result, nil
}

func (rate *ReferenceRate) RateFor(currency string) *float64 {
	for _, r := range rate.Rates {
		if r.Currency == currency {
			return &r.Rate
		}
	}
	return nil
}
