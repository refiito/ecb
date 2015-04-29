package ecb

import (
	"encoding/xml"
	"strconv"
	"time"
  "fmt"

	"github.com/refiito/timeoutclient"
)

const (
	dailyRatesXML     = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
	quarterlyRatesXML = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	allRatesXML       = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml"
)

type (
	envelope struct {
		Data []struct {
			Date  string `xml:"time,attr"`
			Rates []struct {
				Currency string `xml:"currency,attr"`
				Rate     string `xml:"rate,attr"`
			} `xml:"Cube"`
		} `xml:"Cube>Cube"`
	}

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
  xmlURL := allRatesXML
  today := time.Now()

  if isSameDay(date, today) {
    xmlURL = dailyRatesXML
  }
  // 2190 hours is about three months...
  if date.After(today.Add(-2190 * time.Hour)) {
    xmlURL = quarterlyRatesXML
  }
  rates, err := fetchRates(xmlURL)
  if err != nil {
    return nil, err
  }
  
  for _, rate := range rates {
    if isSameDay(rate.Date, date) {
      return &rate, nil
    }
  }
  fmt.Println(date)
  return nil, nil
}

func RateForAt(date time.Time, currency string) (*float64, error) {
  rate, err := RatesAt(date)
  if err != nil {
    return nil, err
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

func fetchXML(xmlURL string) (result envelope, err error) {
	client := timeoutclient.NewTimeoutClient(30*time.Second, 30*time.Second)
	response, err := client.Get(xmlURL)
	if err != nil {
		return
	}
	defer response.Body.Close()
	err = xml.NewDecoder(response.Body).Decode(&result)
	return
}

func fetchRates(xmlURL string) ([]ReferenceRate, error) {
	var result []ReferenceRate

	data, err := fetchXML(xmlURL)
	if err != nil {
		return result, err
	}

	for _, date := range data.Data {
		var refRate ReferenceRate

		refRate.Rates = append(refRate.Rates, Rate{Currency: "EUR", Rate: 1.0})

		if parsedDate, err := time.Parse("2006-01-02", date.Date); err != nil {
			return result, err
		} else {
			refRate.Date = parsedDate
		}

		for _, item := range date.Rates {
			if parsedRate, err := strconv.ParseFloat(item.Rate, 64); err != nil {
				return result, err
			} else {
				refRate.Rates = append(refRate.Rates, Rate{Currency: item.Currency, Rate: parsedRate})
			}
		}
		result = append(result, refRate)
	}
	return result, nil
}
