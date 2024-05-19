package exchange

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	API_URL = "https://api.privatbank.ua/p24api/pubinfo?exchange&coursid=11"
	USD     = "USD"
)

type Rate struct {
	Ccy     string `json:"ccy"`
	BaseCcy string `json:"base_ccy"`
	Buy     string `json:"buy"`
	Sale    string `json:"sale"`
}

var (
	ErrUSDRateNotFound = errors.New("USD rate not found")
)

func FetchRate() (float64, error) {
	resp, err := http.Get(API_URL)
	if err != nil {
		return 0, fmt.Errorf("failed to make GET request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	return parseRates(body)
}

func parseRates(body []byte) (float64, error) {
	var rates []Rate
	err := json.Unmarshal(body, &rates)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return findUSDRate(rates)
}

func findUSDRate(rates []Rate) (float64, error) {
	for _, rate := range rates {
		if rate.Ccy == USD {
			sale, err := strconv.ParseFloat(rate.Sale, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse sale rate: %w", err)
			}
			return sale, nil
		}
	}

	return 0, ErrUSDRateNotFound
}
