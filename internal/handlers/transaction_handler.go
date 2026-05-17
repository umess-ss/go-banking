package handlers

import (
	"go-banking/pkg/response"
	"net/http"
	"strconv"

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
	transactions, err := h.service.GetTransactions(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve transactions")
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "Transactions retrieved successfully", transactions)
}

func (h *TransactionHandler) GetTransactionsByAccountID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	accountID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	transactions, err := h.service.GetTransactionsByAccountID(r.Context(), accountID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve transactions")
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "Transactions retrieved successfully", transactions)
}
