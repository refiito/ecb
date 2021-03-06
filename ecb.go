package ecb

import (
	"errors"
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

func CurrencyRateAt(date time.Time, currency string) (*CurrencyRate, error) {
	if date.Before(time.Date(1999, 1, 4, 0, 0, 0, 0, time.UTC)) {
		return nil, errors.New("Date before data start")
	}

	for i := 0; i < 5; i++ {
		rates, err := RatesAt(date.AddDate(0, 0, -1*i))
		if err != nil {
			return nil, err
		}
		if rates == nil {
			continue
		}

		// Don't output pointer to value
		rateFloat := *rates.RateFor(currency)
		return &CurrencyRate{Date: date, Currency: currency, Rate: &rateFloat}, nil
	}

	return nil, nil
}

func FilledCurrencyRatesBetween(rangeStart, rangeEnd time.Time, currency string) ([]*CurrencyRate, error) {
	checkDate := rangeStart
	daysNum := int(rangeEnd.Sub(rangeStart).Hours() / 24)

	rates := make([]*CurrencyRate, 0, daysNum)
	for i := 0; i < daysNum; i++ {
		rateAt, err := CurrencyRateAt(checkDate.AddDate(0, 0, i), currency)
		if err != nil {
			return nil, err
		}
		rates = append(rates, rateAt)
	}
	return rates, nil
}

func PreWarmCache(rangeStart, rangeEnd time.Time) error {
	return rateCache.populate(rangeStart, rangeEnd)
}

func RateForAt(date time.Time, currency string) (*float64, error) {
	rate, err := RatesAt(date)
	if err != nil || rate == nil {
		return nil, err
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

func FilledRatesForBetween(rangeStart, rangeEnd time.Time, currency string) ([]*CurrencyRate, error) {
	// Warm the cache, not caring about the return here, only the error
	_, err := rateCache.ratesBetween(rangeStart, rangeEnd)
	if err != nil {
		return nil, err
	}

	var result []*CurrencyRate
	for checkDate := rangeStart; !isOneDayLater(checkDate, rangeEnd); checkDate.AddDate(0, 0, 1) {
		rateAt, err := CurrencyRateAt(checkDate, currency)
		if err != nil {
			return nil, err
		}
		result = append(result, rateAt)
	}
	return result, err
}

func (rate *ReferenceRate) RateFor(currency string) *float64 {
	for _, r := range rate.Rates {
		if r.Currency == currency {
			tmp := r.Rate // Don't return a direct pointer
			return &tmp
		}
	}
	return nil
}

func isOneDayLater(t1, t2 time.Time) bool {
	return isSameDay(t1.AddDate(0, 0, 1), t2)
}
