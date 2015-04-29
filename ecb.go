package ecb

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"time"
)

// ECB XML envelope
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

func fetchXML() (result envelope, err error) {
	response, err := http.Get("http://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml")
	if err != nil {
		return
	}
	defer response.Body.Close()
	err = xml.NewDecoder(response.Body).Decode(&result)
	return
}

func FetchRates() ([]ReferenceRate, error) {
	var result []ReferenceRate

	data, err := fetchXML()
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
