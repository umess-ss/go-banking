package handlers

import (
	"encoding/json"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "ok",
		"message": "Go Banking API is healthy",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
