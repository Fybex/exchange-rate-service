package email_sender

import (
	"errors"
	"testing"

	"github.com/Fybex/exchange-rate-service/pkg/exchange"
	"github.com/Fybex/exchange-rate-service/pkg/models"
)

var (
	mockGetSubscribers func() ([]models.Subscriber, error)
	mockFetchRate      func() (float64, error)
	mockSendEmail      func(recipient string, rate float64) error
)

func init() {
	models.GetSubscribers = func() ([]models.Subscriber, error) {
		return mockGetSubscribers()
	}
	exchange.FetchRate = func() (float64, error) {
		return mockFetchRate()
	}
	sendEmail = func(recipient string, rate float64) error {
		return mockSendEmail(recipient, rate)
	}
}

func TestSendRateEmailsSuccess(t *testing.T) {
	mockGetSubscribers = func() ([]models.Subscriber, error) {
		return []models.Subscriber{
			{ID: 1, Email: "test1@example.com"},
			{ID: 2, Email: "test2@example.com"},
		}, nil
	}

	mockFetchRate = func() (float64, error) {
		return 39.84, nil
	}

	mockSendEmail = func(recipient string, rate float64) error {
		return nil
	}

	err := SendRateEmails()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSendRateEmailsFetchSubscribersError(t *testing.T) {
	mockGetSubscribers = func() ([]models.Subscriber, error) {
		return nil, errors.New("database error")
	}

	mockFetchRate = func() (float64, error) {
		return 0, nil
	}

	mockSendEmail = func(recipient string, rate float64) error {
		return nil
	}

	err := SendRateEmails()
	if err == nil || err.Error() != "failed to get subscribers: database error" {
		t.Fatalf("Expected database error, got %v", err)
	}
}

func TestSendRateEmailsFetchRateError(t *testing.T) {
	mockGetSubscribers = func() ([]models.Subscriber, error) {
		return []models.Subscriber{
			{ID: 1, Email: "test1@example.com"},
			{ID: 2, Email: "test2@example.com"},
		}, nil
	}
	mockFetchRate = func() (float64, error) {
		return 0, errors.New("exchange rate error")
	}
	mockSendEmail = func(recipient string, rate float64) error {
		return nil
	}

	err := SendRateEmails()
	if err == nil || err.Error() != "failed to fetch exchange rate: exchange rate error" {
		t.Fatalf("Expected exchange rate error, got %v", err)
	}
}

func TestSendRateEmailsSendEmailError(t *testing.T) {
	mockGetSubscribers = func() ([]models.Subscriber, error) {
		return []models.Subscriber{
			{ID: 1, Email: "test1@example.com"},
			{ID: 2, Email: "test2@example.com"},
		}, nil
	}

	mockFetchRate = func() (float64, error) {
		return 39.84, nil
	}

	emailSent := make(map[string]bool)

	mockSendEmail = func(recipient string, _ float64) error {
		emailSent[recipient] = true
		if recipient == "test1@example.com" {
			return errors.New("failed to send email")
		}
		return nil
	}

	err := SendRateEmails()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !emailSent["test1@example.com"] || !emailSent["test2@example.com"] {
		t.Errorf("Expected emails to be sent, but not all were sent")
	}
}
