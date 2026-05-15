package main

import (
	"fmt"
	"go-banking/internal/handlers"
	"go-banking/internal/repository"
	"go-banking/internal/services"
	"log"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	accountRepo := repository.NewAccountRepository()
	transactionRepo := repository.NewTransactionRepository()
	accountService := services.NewAccountService(accountRepo, transactionRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

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
		path := r.URL.Path

		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(path, "/transactions"):
			transactionHandler.GetTransactionsByAccountID(w, r)

		case r.Method == http.MethodGet:
			accountHandler.GetAccountByID(w, r)

		case r.Method == http.MethodPost && strings.HasSuffix(path, "/deposit"):
			accountHandler.Deposit(w, r)

		case r.Method == http.MethodPost && strings.HasSuffix(path, "/withdraw"):
			accountHandler.Withdraw(w, r)

		default:
			http.Error(w, "route not found", http.StatusNotFound)
		}
	})

	mux.HandleFunc("/transfer", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		} else {
			accountHandler.Transfer(w, r)
		}
	})

	mux.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			transactionHandler.GetTransactions(w, r)
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
