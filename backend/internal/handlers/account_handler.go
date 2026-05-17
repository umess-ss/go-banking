package handlers

import (
	"encoding/json"
	"go-banking/internal/middleware"
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

func getAccountIDFromRoute(r *http.Request) (int64, error) {
	idText := chi.URLParam(r, "id")
	return strconv.ParseInt(idText, 10, 64)
}

func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	accounts, err := h.service.GetAccounts(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to fetch accounts")
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "accounts fetched successfully", accounts)
}

func (h *AccountHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := getAccountIDFromRoute(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	account, err := h.service.GetAccountByID(r.Context(), id, userID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "account not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "account fetched successfully", account)
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var account models.Account

	err = json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createdAccount, err := h.service.CreateAccount(r.Context(), userID, account)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, true, "account created successfully", createdAccount)
}

func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	accountID, err := getAccountIDFromRoute(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	var request models.AmountRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	account, err := h.service.Deposit(r.Context(), userID, accountID, request.Amount)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "deposit completed successfully", account)
}

func (h *AccountHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	accountID, err := getAccountIDFromRoute(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	var request models.AmountRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	account, err := h.service.Withdraw(r.Context(), userID, accountID, request.Amount)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "withdraw completed successfully", account)
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var request models.TransferRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.service.Transfer(r.Context(), userID, request.FromAccountID, request.ToAccountID, request.Amount)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "transfer completed successfully", nil)
}
