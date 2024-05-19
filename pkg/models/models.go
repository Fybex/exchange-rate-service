package models

import (
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type Subscriber struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

var (
	ErrSubscriberExists = errors.New("subscriber already exists")
	ErrUniqueViolation  = "23505" // PostgreSQL unique violation error code
)

func AddSubscriber(subscriber Subscriber) error {
	stmt, err := db.Prepare("INSERT INTO subscribers(email) VALUES($1)")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(subscriber.Email)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == pq.ErrorCode(ErrUniqueViolation) {
				return ErrSubscriberExists
			}
		}
		return fmt.Errorf("could not execute statement: %w", err)
	}
	return nil
}

func GetSubscribers() ([]Subscriber, error) {
	rows, err := db.Query("SELECT id, email FROM subscribers")
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	var subscribers []Subscriber
	for rows.Next() {
		var subscriber Subscriber
		err = rows.Scan(&subscriber.ID, &subscriber.Email)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		subscribers = append(subscribers, subscriber)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over rows: %w", err)
	}

	return subscribers, nil
}
