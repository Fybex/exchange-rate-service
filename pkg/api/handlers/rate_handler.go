package handlers

import (
	"encoding/json"
	"github.com/Fybex/exchange-rate-service/pkg/exchange"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Rate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rate, err := exchange.FetchRate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		log.Printf("Failed fetching exchange rate: %v", err)
		return
	}
	json.NewEncoder(w).Encode(rate)
}
