package main

import (
	"github.com/Fybex/exchange-rate-service/pkg/api/handlers"
	"github.com/Fybex/exchange-rate-service/pkg/email_sender"
	"github.com/Fybex/exchange-rate-service/pkg/models"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	models.InitDB()

	router := mux.NewRouter()
	router.HandleFunc("/api/rate", handlers.Rate).Methods("GET")
	router.HandleFunc("/api/subscribe", handlers.Subscribe).Methods("POST")

	email_sender.ScheduleEmails()

	log.Fatal(http.ListenAndServe(":8000", router))
}
