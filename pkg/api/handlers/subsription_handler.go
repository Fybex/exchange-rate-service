package handlers

import (
	"log"
	"net/http"
	"net/mail"

	"github.com/Fybex/exchange-rate-service/pkg/models"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Missing email field", http.StatusBadRequest)
		log.Println("Missing email field")
		return
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		log.Println("Invalid email format")
		return
	}

	subscriber := models.Subscriber{Email: email}
	err = models.AddSubscriber(subscriber)
	if err != nil {
		if err == models.ErrSubscriberExists {
			http.Error(w, "Email already exists", http.StatusConflict)
			log.Println("Email already exists")
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Failed to add subscriber: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("E-mail added"))
}
