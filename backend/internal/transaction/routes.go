package transaction

import (
	"go-banking/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *TransactionHandler) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/transactions", handler.GetTransactions)
	})
}
