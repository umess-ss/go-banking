package main

import (
	"fmt"
	"go-banking/internal/account"
	"go-banking/internal/auth"
	"go-banking/internal/config"
	"go-banking/internal/database"
	"go-banking/internal/health"
	"go-banking/internal/middleware"
	"go-banking/internal/response"
	"go-banking/internal/transaction"
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

	accountRepo := account.NewAccountRepository(dbPool)
	transactionRepo := transaction.NewTransactionRepository(dbPool)
	userRepo := auth.NewUserRepository(dbPool)

	accountService := account.NewAccountService(accountRepo, transactionRepo)
	transactionService := transaction.NewTransactionService(transactionRepo)
	authService := auth.NewAuthService(userRepo)

	accountHandler := account.NewAccountHandler(accountService)
	transactionHandler := transaction.NewTransactionHandler(transactionService)
	authHandler := auth.NewAuthHandler(authService)

	health.RegisterRoutes(router)
	auth.RegisterRoutes(router, authHandler)
	account.RegisterRoutes(router, accountHandler, transactionHandler)
	transaction.RegisterRoutes(router, transactionHandler)

	port := ":" + cfg.Port

	fmt.Println("Go Banking API started on port", port)

	handler := middleware.CORS(router)
	err = http.ListenAndServe(port, handler)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
