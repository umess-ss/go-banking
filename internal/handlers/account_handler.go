package handlers

import (
	"encoding/json"
	"go-banking/internal/models"
	"go-banking/internal/services"
	"go-banking/pkg/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// handler for account-related HTTP requests. It uses the AccountService to perform business
// logic and interacts with the HTTP request and response to provide the appropriate outtput.
// It includes methods for handling requests to get all accounts, get an account by ID, and
// create a new account. Each method decodes the request, calls the service layer, and encodes
// the response as JSON. It also handles error cases and returns apropriate HTTP status codes and messages.
type AccountHandler struct {
	service *services.AccountService
}

func NewAccountHandler(service *services.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts := h.service.GetAllAccounts()
	response.WriteJSON(w, http.StatusOK, true, "Accounts retrieved successfully", accounts)
}

func (h *AccountHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	id, err := getAccountIDFromRouter(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	account, err := h.service.GetAccountByID(id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "account not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "Account retrieved successfully", account)
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account models.Account

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createdAccount, err := h.service.CreateAccount(account)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, true, "Account created successfully", createdAccount)
}

func getAccountIDFromRouter(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	return strconv.Atoi(idStr)
}

func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	id, err := getAccountIDFromRouter(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	var request struct {
		Amount float64 `json:"amount"`
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	account, err := h.service.Deposit(id, request.Amount)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "Deposit successful", account)
}

func (h *AccountHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	id, err := getAccountIDFromRouter(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	var request struct {
		Amount float64 `json:"amount"`
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	account, err := h.service.Withdraw(id, request.Amount)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "Withdrawal successful", account)
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var request struct {
		FromAccountID int     `json:"from_account_id"`
		ToAccountID   int     `json:"to_account_id"`
		Amount        float64 `json:"amount"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.service.Transfer(request.FromAccountID, request.ToAccountID, request.Amount)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "Transfer successful", nil)
}
