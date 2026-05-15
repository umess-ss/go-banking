package main

import (
	"fmt"
	"go-banking/internal/handlers"
	"go-banking/internal/repository"
	"go-banking/internal/services"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	accountRepo := repository.NewAccountRepository()
	accountService := services.NewAccountService(accountRepo)
	accountHandler := handlers.NewAccountHandler(accountService)

	mux.HandleFunc("/health", handlers.HealthCheckHandler)

	mux.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			accountHandler.GetAllAccounts(w, r)
		} else if r.Method == http.MethodPost {
			accountHandler.CreateAccount(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/accounts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			accountHandler.GetAccountByID(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := 8080
	fmt.Printf("Starting Go Banking API on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
