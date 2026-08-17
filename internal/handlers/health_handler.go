package handlers

import (
	"net/http"
	"subscriptions-api-postgres/internal/response"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	result := response.ResponseSuccess{
		Message: "subscriptions-api-postgres is running",
	}

	response.WriteJSON(w, http.StatusOK, result)
}
