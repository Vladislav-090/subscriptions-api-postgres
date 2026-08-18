package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseSuccess struct {
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func WriteError(w http.ResponseWriter, status int, message string) {
	errorResponse := ResponseError{
		Error: message,
	}

	WriteJSON(w, status, errorResponse)
}

func WriteServerError(w http.ResponseWriter, message string, err error) {
	slog.Error(message, "error", err)
	WriteError(w, http.StatusInternalServerError, message)
}
