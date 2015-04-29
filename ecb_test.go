package ecb

import (
	"testing"
)

func TestFetchRates(t *testing.T) {
	_, err := FetchRates()
	if err != nil {
		t.Error(err)
	}
}
