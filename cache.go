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
	rates := make(chan *ReferenceRate)
	err := fetchRates(xmlURL, rates)
	if err != nil {
		return err
	}
	for {
		rate := <-rates
		if rate == nil {
			break
		}
		if flow != nil {
			flow <- rate
		}
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
	rates := make(chan *ReferenceRate, 1)
	err := cache.fetch(date, date, rates)
	if len(rates) == 0 {
		return nil, err
	}
	return <-rates, err
}

func (cache *cachedRates) fetchWithPopulate(rangeStart, rangeEnd time.Time, result chan *ReferenceRate) (err error) {
	rates := make(chan *ReferenceRate)
	err = cache.populate(rangeStart, rangeEnd, rates)
	for {
		rate := <-rates
		if rate == nil {
			break
		}
		if (rate.Date.After(rangeStart) && rate.Date.Before(rangeEnd)) || isSameDay(rangeStart, rate.Date) || isSameDay(rangeEnd, rate.Date) {
			result <- rate
		}
	}
	return
}

func (cache *cachedRates) rangeCached(rangeStart, rangeEnd time.Time) bool {
	return cache.rates != nil && !cache.start.IsZero() && !cache.end.IsZero() && !rangeStart.Before(cache.start) && !rangeEnd.After(cache.end)
}

func (cache *cachedRates) fetchFromCache(rangeStart, rangeEnd time.Time, result chan *ReferenceRate) (err error) {
	if !cache.rangeCached(rangeStart, rangeEnd) {
		return cache.fetchWithPopulate(rangeStart, rangeEnd, result)
	}
	for _, rate := range cache.rates {
		if (rate.Date.After(rangeStart) && rate.Date.Before(rangeEnd)) || isSameDay(rangeStart, rate.Date) || isSameDay(rangeEnd, rate.Date) {
			result <- rate
		}
	}
	return
}

func (cache *cachedRates) fetch(rangeStart, rangeEnd time.Time, result chan *ReferenceRate) (err error) {
	if !cache.rangeCached(rangeStart, rangeEnd) {
		return cache.fetchWithPopulate(rangeStart, rangeEnd, result)
	}
	return cache.fetchFromCache(rangeStart, rangeEnd, result)
}

func (cache *cachedRates) ratesBetween(rangeStart, rangeEnd time.Time) (result []ReferenceRate, err error) {
	if cache.start.IsZero() || cache.end.IsZero() || rangeStart.Before(cache.start) || rangeEnd.After(cache.end) {
		err = cache.populate(rangeStart, rangeEnd, nil)
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

func isSameDay(t1, t2 time.Time) (result bool) {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
