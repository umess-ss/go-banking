package handlers

import (
	"encoding/json"
	"go-banking/internal/middleware"
	"net/http"

	"go-banking/internal/models"
	"go-banking/internal/services"
	"go-banking/pkg/response"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request models.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Register(r.Context(), request)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, true, "user registered successfully", user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.service.Login(r.Context(), request)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "login successful", result)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.service.GetCurrentUser(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "current user fetched successfully", user)
}
