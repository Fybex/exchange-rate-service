package exchange

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

const mockAPIResponse = `[{"ccy":"USD","base_ccy":"UAH","buy":"39.25000","sale":"39.84064"},{"ccy":"EUR","base_ccy":"UAH","buy":"42.52000","sale":"43.29004"}]`

const expectedRate = 39.84064

func TestFetchRate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockAPIResponse))
	}))
	defer server.Close()

	originalAPIURL := API_URL
	API_URL = server.URL
	defer func() { API_URL = originalAPIURL }()

	rate, err := FetchRate()
	if err != nil {
		t.Fatalf("FetchRate() returned an error: %v", err)
	}

	if rate != expectedRate {
		t.Errorf("FetchRate() returned %v, expected %v", rate, expectedRate)
	}
}

func TestParseRates(t *testing.T) {
	rate, err := parseRates([]byte(mockAPIResponse))
	if err != nil {
		t.Fatalf("parseRates() returned an error: %v", err)
	}

	if rate != expectedRate {
		t.Errorf("parseRates() returned %v, expected %v", rate, expectedRate)
	}
}

func TestFindUSDRate(t *testing.T) {
	rates := []Rate{
		{Ccy: "USD", BaseCcy: "UAH", Buy: "39.25000", Sale: "39.84064"},
		{Ccy: "EUR", BaseCcy: "UAH", Buy: "42.52000", Sale: "43.29004"},
	}

	rate, err := findUSDRate(rates)
	if err != nil {
		t.Fatalf("findUSDRate() returned an error: %v", err)
	}

	if rate != expectedRate {
		t.Errorf("findUSDRate() returned %v, expected %v", rate, expectedRate)
	}
}

func TestFindUSDRateNotFound(t *testing.T) {
	rates := []Rate{
		{Ccy: "EUR", BaseCcy: "UAH", Buy: "42.52000", Sale: "43.29004"},
	}

	_, err := findUSDRate(rates)
	if !errors.Is(err, ErrUSDRateNotFound) {
		t.Errorf("findUSDRate() returned %v, expected %v", err, ErrUSDRateNotFound)
	}
}
