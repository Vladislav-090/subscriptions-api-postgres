package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"subscriptions-api-postgres/internal/auth"
	"subscriptions-api-postgres/internal/models"
	"subscriptions-api-postgres/internal/response"
	"subscriptions-api-postgres/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input models.UserRegisterInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	user, err := h.authService.Register(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailRequired):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrPasswordTooShort):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrEmailTaken):
			response.WriteError(w, http.StatusConflict, err.Error())
		default:
			response.WriteServerError(w, "failed to register user", err)
		}
		return
	}

	response.WriteJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.UserLoginInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	token, err := h.authService.Login(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			response.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			response.WriteServerError(w, "failed to login", err)
		}
		return
	}

	response.WriteJSON(w, http.StatusOK, loginResponse{Token: token})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, errMissingClaims.Error())
		return
	}

	if err := h.authService.Logout(r.Context(), claims); err != nil {
		response.WriteServerError(w, "failed to logout", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
