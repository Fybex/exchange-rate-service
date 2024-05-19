package models

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestAddSubscriber(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %v", err)
	}
	defer mockDB.Close()
	SetDB(mockDB)

	tests := []struct {
		name          string
		subscriber    Subscriber
		mockExpect    func(mock sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name:       "AddSubscriber success",
			subscriber: Subscriber{Email: "test@example.com"},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO subscribers").
					ExpectExec().
					WithArgs("test@example.com").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:       "AddSubscriber unique violation",
			subscriber: Subscriber{Email: "test@example.com"},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO subscribers").
					ExpectExec().
					WithArgs("test@example.com").
					WillReturnError(&pq.Error{Code: pq.ErrorCode(ErrUniqueViolation)})
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect(mock)
			err := AddSubscriber(tt.subscriber)
			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err != nil)
			}
		})
	}
}

func TestGetSubscribers(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %v", err)
	}
	defer mockDB.Close()
	SetDB(mockDB)

	queryError := errors.New("query error")

	tests := []struct {
		name           string
		mockExpect     func(mock sqlmock.Sqlmock)
		expectedResult []Subscriber
		expectedError  bool
	}{
		{
			name: "GetSubscribers success",
			mockExpect: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email"}).
					AddRow(1, "test1@example.com").
					AddRow(2, "test2@example.com")
				mock.ExpectQuery("SELECT id, email FROM subscribers").
					WillReturnRows(rows)
			},
			expectedResult: []Subscriber{
				{ID: 1, Email: "test1@example.com"},
				{ID: 2, Email: "test2@example.com"},
			},
			expectedError: false,
		},
		{
			name: "GetSubscribers query error",
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, email FROM subscribers").
					WillReturnError(queryError)
			},
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect(mock)
			result, err := GetSubscribers()
			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err != nil)
			}
			if len(result) != len(tt.expectedResult) {
				t.Errorf("Expected result length %v, got %v", len(tt.expectedResult), len(result))
			}
			for i, subscriber := range result {
				if subscriber != tt.expectedResult[i] {
					t.Errorf("Expected result %v, got %v", tt.expectedResult[i], subscriber)
				}
			}
		})
	}
}
