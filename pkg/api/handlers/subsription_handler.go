package handlers

import (
	"exchange-rate-service/pkg/models"
	"net/http"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Missing email field", http.StatusBadRequest)
		return
	}

	subscriber := models.Subscriber{Email: email}
	err := models.AddSubscriber(subscriber)
	if err != nil {
		if err == models.ErrSubscriberExists {
			http.Error(w, "Email already exists", http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("E-mail added"))
}
