package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Fybex/exchange-rate-service/pkg/models"
)

var (
	mockAddSubscriber func(subscriber models.Subscriber) error
)

func init() {
	models.AddSubscriber = func(subscriber models.Subscriber) error {
		return mockAddSubscriber(subscriber)
	}
}

func TestSubscribeSuccess(t *testing.T) {
	mockAddSubscriber = func(subscriber models.Subscriber) error {
		return nil
	}

	req, err := http.NewRequest("POST", "/subscribe", strings.NewReader("email=test@example.com"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Subscribe)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "E-mail added"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestSubscribeMissingEmail(t *testing.T) {
	req, err := http.NewRequest("POST", "/subscribe", strings.NewReader(""))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Subscribe)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := "Missing email field\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestSubscribeEmailExists(t *testing.T) {
	mockAddSubscriber = func(subscriber models.Subscriber) error {
		return models.ErrSubscriberExists
	}

	req, err := http.NewRequest("POST", "/subscribe", strings.NewReader("email=test@example.com"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Subscribe)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}

	expected := "Email already exists\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestSubscribeInternalError(t *testing.T) {
	mockAddSubscriber = func(_ models.Subscriber) error {
		return errors.New("internal error")
	}

	req, err := http.NewRequest("POST", "/subscribe", strings.NewReader("email=test@example.com"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Subscribe)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := "Internal server error\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
