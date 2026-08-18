package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"subscriptions-api-postgres/internal/auth"
	"subscriptions-api-postgres/internal/models"
	"subscriptions-api-postgres/internal/response"
	"subscriptions-api-postgres/internal/service"
)

var errMissingClaims = errors.New("missing auth claims")

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

	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, errMissingClaims.Error())
		return
	}

	if claims.Role != models.RoleAdmin {
		input.UserID = claims.UserID
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
			response.WriteServerError(w, "failed to create subscription", err)
		}
		return
	}

	response.WriteJSON(w, http.StatusCreated, createdSubscription)
}

func (h *SubscriptionHandler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {

	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, errMissingClaims.Error())
		return
	}

	userID := claims.UserID
	if claims.Role == models.RoleAdmin {
		parsedID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid user ID")
			return
		}
		userID = parsedID
	}

	subscriptions, err := h.subscriptionService.GetSubscriptions(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUserID):
			response.WriteError(w, http.StatusBadRequest, "invalid user ID")
		default:
			response.WriteServerError(w, "failed to get subscriptions", err)
		}
		return
	}

	response.WriteJSON(w, http.StatusOK, subscriptions)
}


func (h *SubscriptionHandler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}

	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, errMissingClaims.Error())
		return
	}

	userID := claims.UserID
	if claims.Role == models.RoleAdmin {
		parsedID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid user ID")
			return
		}
		userID = parsedID
	}

	subscription, err := h.subscriptionService.GetSubscriptionByID(r.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSubscriptionID):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrInvalidUserID):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, sql.ErrNoRows):
			response.WriteError(w, http.StatusNotFound, "subscription not found")
		default:
			response.WriteServerError(w, "failed to get subscription", err)
		}
		return
	}

	response.WriteJSON(w, http.StatusOK, subscription)
}

func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var input models.SubscriptionUpdateInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	updatedSubscription, err := h.subscriptionService.UpdateSubscription(r.Context(), id, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSubscriptionID):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrInvalidUserID):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrInvalidPrice):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrEmptyService):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, sql.ErrNoRows):
			response.WriteError(w, http.StatusNotFound, "subscription not found")
		default:
			response.WriteServerError(w, "failed to update subscription", err)
		}
		return
	}

	response.WriteJSON(w, http.StatusOK, updatedSubscription)
}

func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	err = h.subscriptionService.DeleteSubscription(r.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSubscriptionID):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrInvalidUserID):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, sql.ErrNoRows):
			response.WriteError(w, http.StatusNotFound, "subscription not found")
		default:
			response.WriteServerError(w, "failed to delete subscription", err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
