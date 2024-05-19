package main

import (
	"exchange-rate-service/pkg/api/handlers"
	"exchange-rate-service/pkg/email_sender"
	"exchange-rate-service/pkg/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	models.InitDB()

	router := mux.NewRouter()
	router.HandleFunc("/api/rate", handlers.GetExchangeRate).Methods("GET")
	router.HandleFunc("/api/subscribe", handlers.Subscribe).Methods("POST")

	email_sender.ScheduleEmails()

	log.Fatal(http.ListenAndServe(":8000", router))
}
