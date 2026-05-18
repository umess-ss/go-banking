package health

import (
	"go-banking/internal/response"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
	}
	response.WriteJSON(w, http.StatusOK, true, "Go Banking API is healthy", data)
}
