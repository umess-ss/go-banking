package health

import (
	"context"
	"net/http"
	"time"

	"go-banking/internal/response"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{
		db: db,
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, true, "service is healthy", map[string]string{
		"status":  "ok",
		"service": "go-banking-api",
	})
}

func (h *Handler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		response.WriteError(w, http.StatusServiceUnavailable, "database not ready")
		return
	}

	response.WriteJSON(w, http.StatusOK, true, "service is ready", map[string]string{
		"status":   "ready",
		"database": "connected",
	})
}
