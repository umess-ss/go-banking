package main

import (
	"fmt"
	"go-banking/internal/handlers"
	"go-banking/internal/repository"
	"go-banking/internal/services"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	router := chi.NewRouter()

	accountRepo := repository.NewAccountRepository()
	transactionRepo := repository.NewTransactionRepository()
	accountService := services.NewAccountService(accountRepo, transactionRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	router.Get("/health", handlers.HealthCheckHandler)

	router.Route("/accounts", func(r chi.Router) {
		r.Get("/", accountHandler.GetAllAccounts)
		r.Post("/", accountHandler.CreateAccount)

		r.Get("/{id}", accountHandler.GetAccountByID)
		r.Post("/{id}/deposit", accountHandler.Deposit)
		r.Post("/{id}/withdraw", accountHandler.Withdraw)
		r.Get("/{id}/transactions", transactionHandler.GetTransactionsByAccountID)
	})

	router.Post("/transfer", accountHandler.Transfer)

	router.Get("/transactions", transactionHandler.GetTransactions)

	port := ":8080"

	fmt.Println("Go Banking API started on port", port)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
