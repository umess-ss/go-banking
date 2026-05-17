package main

import (
	"fmt"
	"go-banking/internal/config"
	"go-banking/internal/database"
	"go-banking/internal/handlers"
	"go-banking/internal/middleware"
	"go-banking/internal/repository"
	"go-banking/internal/services"
	"go-banking/pkg/response"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.Load()

	dbPool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	fmt.Println("Connected to database successfully")

	router := chi.NewRouter()

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.WriteError(w, http.StatusNotFound, "route not found")
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	})

	router.Use(middleware.Logger)
	router.Use(middleware.Recovery)

	accountRepo := repository.NewAccountRepository(dbPool)
	transactionRepo := repository.NewTransactionRepository(dbPool)
	userRepo := repository.NewUserRepository(dbPool)

	accountService := services.NewAccountService(accountRepo, transactionRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	authService := services.NewAuthService(userRepo)

	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	authHandler := handlers.NewAuthHandler(authService)

	router.Get("/health", handlers.HealthCheckHandler)

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Route("/accounts", func(r chi.Router) {
			r.Get("/", accountHandler.GetAllAccounts)
			r.Post("/", accountHandler.CreateAccount)

			r.Get("/{id}", accountHandler.GetAccountByID)
			r.Post("/{id}/deposit", accountHandler.Deposit)
			r.Post("/{id}/withdraw", accountHandler.Withdraw)
			r.Get("/{id}/transactions", transactionHandler.GetTransactionsByAccountID)
		})

		r.Post("/transfer", accountHandler.Transfer)
		r.Get("/transactions", transactionHandler.GetTransactions)

	})

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	port := ":" + cfg.Port

	fmt.Println("Go Banking API started on port", port)

	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
