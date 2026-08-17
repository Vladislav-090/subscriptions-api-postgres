package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"subscriptions-api-postgres/internal/models"
	"subscriptions-api-postgres/internal/response"
	"subscriptions-api-postgres/internal/service"
)

type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
}

func NewSubscriptionsHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: service,
	}
}

func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var input models.SubscriptionInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	createdSubscription, err := h.subscriptionService.CreateSubscription(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmptyService):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrInvalidPrice):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrInvalidUserID):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			response.WriteError(w, http.StatusInternalServerError, "failed to create subscription")
		}
		return
	}

	response.WriteJSON(w, http.StatusCreated, createdSubscription)
}

func (h *SubscriptionHandler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	subscriptions, err := h.subscriptionService.GetSubscriptions(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUserID):
			response.WriteError(w, http.StatusBadRequest, "invalid user ID")
		default:
			response.WriteError(w, http.StatusInternalServerError, "failed to get subscriptions")
		}
		return
	}

	response.WriteJSON(w, http.StatusOK, subscriptions)
}


func (h *SubscriptionHandler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request){

}
