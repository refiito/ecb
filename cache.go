package ecb

import (
	"time"
)

type cachedRates struct {
	start time.Time
	end   time.Time
	rates []ReferenceRate
}

var rateCache cachedRates

func (cache *cachedRates) populate(rangeStart, rangeEnd time.Time) error {
	xmlURL := allRatesXML
	today := time.Now()

	if isSameDay(rangeStart, rangeEnd) && isSameDay(rangeStart, today) {
		xmlURL = dailyRatesXML
	}
	// 2190 hours is about three months...
	if rangeStart.After(today.Add(-2190 * time.Hour)) {
		xmlURL = quarterlyRatesXML
	}
	rates, err := fetchRates(xmlURL)
	if err != nil {
		return err
	}
	cache.rates = rates
	if len(cache.rates) > 0 {
		cache.start = cache.rates[len(cache.rates)-1].Date
		cache.end = cache.rates[0].Date
	}

	return nil
}

func (cache *cachedRates) ratesAt(date time.Time) (*ReferenceRate, error) {
	rates, err := cache.ratesBetween(date, date)
	if len(rates) == 0 {
		return nil, err
	}
	return &rates[0], err
}

func (cache *cachedRates) ratesBetween(rangeStart, rangeEnd time.Time) (result []ReferenceRate, err error) {
	if cache.start.IsZero() || cache.end.IsZero() || rangeStart.Before(cache.start) || rangeEnd.After(cache.end) {
		err = cache.populate(rangeStart, rangeEnd)
		if err != nil {
			return
		}
	}
	for _, rate := range cache.rates {
		if (rate.Date.After(rangeStart) && rate.Date.Before(rangeEnd)) || isSameDay(rangeStart, rate.Date) || isSameDay(rangeEnd, rate.Date) {
			result = append(result, rate)
		}
	}
	return
}

func isSameDay(t1, t2 time.Time) (result bool) {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
