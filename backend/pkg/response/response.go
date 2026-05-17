package response

import (
	"encoding/json"
	"net/http"
)

//It is a utility package that provides functions to standardize the structure of API responses
//and error handling across the application. It defines a struct APIResponse to represent
//the format of the JSON response, and provides two functions: WriteJSON to write a successful
//response and WriteError to write an error response. This helps in maintaining consistency in
//the API responses and simplifies the process of sending responses from different handlers
//in the application.

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(APIResponse{
		Success: statusCode >= 200 && statusCode < 300,
		Message: message,
		Data:    data,
	})
}

func WriteError(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   errorMessage,
	})
}
