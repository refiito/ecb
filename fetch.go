package ecb

import (
	"encoding/xml"
	"strconv"
	"time"

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
)

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

func fetchRates(xmlURL string) ([]*ReferenceRate, error) {
	data, err := fetchXML(xmlURL)
	if err != nil {
		return nil, err
	}

	refRates := make([]*ReferenceRate, 0, len(data.Data))
	for _, date := range data.Data {
		parsedDate, err := time.Parse("2006-01-02", date.Date)
		if err != nil {
			return nil, err
		}

		refRate := ReferenceRate{
			Date: parsedDate.UTC(),
			Rates: []Rate{
				{Currency: "EUR", Rate: 1.0},
			},
		}
		for _, item := range date.Rates {
			if parsedRate, err := strconv.ParseFloat(item.Rate, 64); err != nil {
				return nil, err
			} else {
				refRate.Rates = append(refRate.Rates, Rate{Currency: item.Currency, Rate: parsedRate})
			}
		}

		refRates = append(refRates, &refRate)
	}

	return refRates, nil
}
