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
	transactions := h.service.GetTransactions()

	response.WriteJSON(w, http.StatusOK, true, "Transactions retrieved successfully", transactions)
}

func (h *TransactionHandler) GetTransactionsByAccountID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	accountID, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	transactions := h.service.GetTransactionsByAccountID(accountID)

	response.WriteJSON(w, http.StatusOK, true, "Transactions retrieved successfully", transactions)
}
