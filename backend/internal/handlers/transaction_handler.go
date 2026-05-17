package handlers

import (
	"go-banking/pkg/response"
	"net/http"
	"strconv"

	"go-banking/internal/middleware"
	"go-banking/internal/services"

	"github.com/go-chi/chi/v5"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	transactions, err := h.service.GetTransactions(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to fetch transactions")
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "transactions fetched successfully", transactions)
}

func (h *TransactionHandler) GetTransactionsByAccountID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	idText := chi.URLParam(r, "id")

	accountID, err := strconv.ParseInt(idText, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	transactions, err := h.service.GetTransactionsByAccountID(r.Context(), accountID, userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to fetch account transactions")
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "account transactions fetched successfully", transactions)
}
