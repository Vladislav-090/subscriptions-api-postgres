package router

import (
	"net/http"
	"subscriptions-api-postgres/internal/auth"
	"subscriptions-api-postgres/internal/handlers"
	"subscriptions-api-postgres/internal/middleware"
	"subscriptions-api-postgres/internal/models"
)

func New(
	subscriptionHandler *handlers.SubscriptionHandler,
	authHandler *handlers.AuthHandler,
	jwtManager *auth.JWTManager,
	revoker middleware.TokenRevoker,
) http.Handler {
	mux := http.NewServeMux()

	authenticate := middleware.Authenticate(jwtManager, revoker)
	requireAdmin := middleware.RequireRole(models.RoleAdmin)

	mux.HandleFunc("GET /health", handlers.HealthCheck)

	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.Handle("POST /logout", middleware.Chain(http.HandlerFunc(authHandler.Logout), authenticate))

	mux.Handle("POST /subscriptions", middleware.Chain(http.HandlerFunc(subscriptionHandler.CreateSubscription), authenticate))
	mux.Handle("GET /subscriptions", middleware.Chain(http.HandlerFunc(subscriptionHandler.GetSubscriptions), authenticate))
	mux.Handle("GET /subscriptions/{id}", middleware.Chain(http.HandlerFunc(subscriptionHandler.GetSubscriptionByID), authenticate))
	mux.Handle("PUT /subscriptions/{id}", middleware.Chain(http.HandlerFunc(subscriptionHandler.UpdateSubscription), authenticate, requireAdmin))
	mux.Handle("DELETE /subscriptions/{id}", middleware.Chain(http.HandlerFunc(subscriptionHandler.DeleteSubscription), authenticate, requireAdmin))

	return middleware.Logging(mux)
}
