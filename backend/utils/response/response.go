package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   APIError    `json:"error"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Success(w http.ResponseWriter, r *http.Request, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, r *http.Request, message string, e_message string, code int, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  "error",
		Message: message,
		Error: APIError{
			Code:    code,
			Message: e_message,
		},
	})
}
