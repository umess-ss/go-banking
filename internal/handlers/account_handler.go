package handlers

import (
	"encoding/json"
	"go-banking/internal/models"
	"go-banking/internal/services"
	"net/http"
	"strconv"
	"strings"
)

type AccountHandler struct {
	service *services.AccountService
}

func NewAccountHandler(service *services.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts := h.service.GetAllAccounts()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func (h *AccountHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimPrefix(r.URL.Path, "/accounts/")

	id, err := strconv.Atoi(idText)
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	account, err := h.service.GetAccountByID(id)
	if err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account models.Account

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	createdAccount, err := h.service.CreateAccount(account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdAccount)
}
