package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fybex/exchange-rate-service/pkg/exchange"
)

func mockFetchRateSuccess() (float64, error) {
	return 39.84, nil
}

func mockFetchRateError() (float64, error) {
	return 0, exchange.ErrUSDRateNotFound
}

func TestGetExchangeRateSuccess(t *testing.T) {
	originalFetchRate := exchange.FetchRate
	exchange.FetchRate = mockFetchRateSuccess
	defer func() { exchange.FetchRate = originalFetchRate }()

	req, err := http.NewRequest("GET", "/rate", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Rate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var rate float64
	err = json.NewDecoder(rr.Body).Decode(&rate)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	expectedRate := 39.84
	if rate != expectedRate {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rate, expectedRate)
	}
}

func TestGetExchangeRateError(t *testing.T) {
	originalFetchRate := exchange.FetchRate
	exchange.FetchRate = mockFetchRateError
	defer func() { exchange.FetchRate = originalFetchRate }()

	req, err := http.NewRequest("GET", "/rate", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Rate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	var errorResponse ErrorResponse
	err = json.NewDecoder(rr.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	expectedError := exchange.ErrUSDRateNotFound.Error()
	if errorResponse.Error != expectedError {
		t.Errorf("Handler returned unexpected error: got %v want %v",
			errorResponse.Error, expectedError)
	}
}
