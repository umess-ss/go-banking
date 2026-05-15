package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-banking/internal/services"
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetTransactionsByAccountID(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimPrefix(r.URL.Path, "/accounts/")
	parts := strings.Split(idText, "/")

	accountID, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	transactions := h.service.GetTransactionsByAccountID(accountID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
