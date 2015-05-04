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

func CurrencyRateAt(date time.Time, currency string) (result CurrencyRate, err error) {
	if date.Before(time.Date(1999, 1, 4, 0, 0, 0, 0, time.UTC)) {
		err = errors.New("Date before data start")
		return
	}
	rates, err := RatesAt(date)
	if err != nil {
		return
	}
	if rates == nil {
		return CurrencyRateAt(date.AddDate(0, 0, -1), currency)
	}
	result = CurrencyRate{Date: date, Currency: currency, Rate: rates.RateFor(currency)}
	return
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

func FilledRatesForBetween(rangeStart, rangeEnd time.Time, currency string) ([]CurrencyRate, error) {
	var result []CurrencyRate
	// Warm the cache, not caring about the return here, only the error
	_, err := rateCache.ratesBetween(rangeStart, rangeEnd)
	if err != nil {
		return result, err
	}
	for checkDate := rangeStart; !isOneDayLater(checkDate, rangeEnd); checkDate.AddDate(0, 0, 1) {
		rateAt, err := CurrencyRateAt(checkDate, currency)
		if err != nil {
			return result, err
		}
		result = append(result, rateAt)
	}
	return result, err
}

func (rate *ReferenceRate) RateFor(currency string) *float64 {
	for _, r := range rate.Rates {
		if r.Currency == currency {
			return &r.Rate
		}
	}
	return nil
}

func isOneDayLater(t1, t2 time.Time) bool {
	return isSameDay(t1.AddDate(0, 0, 1), t2)
}
