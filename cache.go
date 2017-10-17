package ecb

import (
	"time"
)

type cachedRates struct {
	start time.Time
	end   time.Time
	rates []*ReferenceRate
}

var rateCache cachedRates

func (cache *cachedRates) populate(rangeStart, rangeEnd time.Time, flow chan *ReferenceRate) error {
	xmlURL := allRatesXML
	today := time.Now()

	if isSameDay(rangeStart, rangeEnd) && isSameDay(rangeStart, today) {
		xmlURL = dailyRatesXML
	}
	// 2190 hours is about three months...
	if rangeStart.After(today.Add(-2190 * time.Hour)) {
		xmlURL = quarterlyRatesXML
	}

	refRates, err := fetchRates(xmlURL)
	if err != nil {
		return err
	}
	for _, rate := range refRates {
		cache.rates = append(cache.rates, rate)
		if cache.start.IsZero() || rate.Date.Before(cache.start) {
			cache.start = rate.Date
		}
		if cache.end.IsZero() || rate.Date.After(cache.end) {
			cache.end = rate.Date
		}
	}

	return nil
}

func (cache *cachedRates) ratesAt(date time.Time) (*ReferenceRate, error) {
	refRates, err := cache.ratesBetween(date, date)
	if err != nil || len(refRates) == 0 {
		return nil, err
	}
	return &refRates[0], nil
}

func (cache *cachedRates) rangeCached(rangeStart, rangeEnd time.Time) bool {
	return cache.rates != nil && !cache.start.IsZero() && !cache.end.IsZero() && !rangeStart.Before(cache.start) && !rangeEnd.After(cache.end)
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
			result = append(result, *rate)
		}
	}
	return
}

func isSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
