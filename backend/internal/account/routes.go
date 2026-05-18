package account

import (
	"go-banking/internal/middleware"
	"go-banking/internal/transaction"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, accountHandler *AccountHandler, transactionHandler *transaction.TransactionHandler) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Route("/accounts", func(r chi.Router) {
			r.Get("/", accountHandler.GetAccounts)
			r.Post("/", accountHandler.CreateAccount)

			r.Get("/{id}", accountHandler.GetAccountByID)
			r.Post("/{id}/deposit", accountHandler.Deposit)
			r.Post("/{id}/withdraw", accountHandler.Withdraw)
			r.Get("/{id}/transactions", transactionHandler.GetTransactionsByAccountID)
		})

		r.Post("/transfer", accountHandler.Transfer)
	})
}
