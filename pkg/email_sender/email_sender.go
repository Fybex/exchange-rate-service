package email_sender

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Fybex/exchange-rate-service/pkg/exchange"
	"github.com/Fybex/exchange-rate-service/pkg/models"

	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

func ScheduleEmails() {
	c := cron.New()
	_, err := c.AddFunc("0 8 * * *", func() { SendRateEmails() })
	if err != nil {
		log.Fatal("Could not schedule email sender: ", err)
	} else {
		log.Println("Email sender scheduled")
	}
	c.Start()
}

func SendRateEmails() error {
	subscribers, err := models.GetSubscribers()
	if err != nil {
		return fmt.Errorf("failed to get subscribers: %w", err)
	}

	exchangeRate, err := exchange.FetchRate()
	if err != nil {
		return fmt.Errorf("failed to fetch exchange rate: %w", err)
	}

	for _, subscriber := range subscribers {
		if err := sendEmail(subscriber.Email, exchangeRate); err != nil {
			log.Printf("Failed to send email to %s: %v", subscriber.Email, err)
		} else {
			log.Printf("Sent email to: %s", subscriber.Email)
		}
	}
	return nil
}

var sendEmail = func(recipient string, rate float64) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "Daily USD to UAH Exchange Rate Update")
	m.SetBody("text/plain", fmt.Sprintf("The current exchange rate from USD to UAH is: %.2f", rate))

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("Failed to parse SMTP_PORT: %v", err)
	}

	d := gomail.NewDialer(os.Getenv("SMTP_SERVER"),
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
