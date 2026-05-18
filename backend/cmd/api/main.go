package main

import (
	"log/slog"
	"net/http"

	"go-banking/internal/account"
	"go-banking/internal/auth"
	"go-banking/internal/config"
	"go-banking/internal/database"
	"go-banking/internal/health"
	"go-banking/internal/logger"
	"go-banking/internal/middleware"
	"go-banking/internal/response"
	"go-banking/internal/transaction"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.AppEnv)

	dbPool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		return
	}
	defer dbPool.Close()

	log.Info("connected to database successfully")

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

	log.Info(
		"server starting",
		slog.String("port", cfg.Port),
		slog.String("env", cfg.AppEnv),
	)

	handler := middleware.CORS(router)

	err = http.ListenAndServe(port, handler)
	if err != nil {
		log.Error("failed to start server", slog.String("error", err.Error()))
	}
}
