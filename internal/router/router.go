package router

import (
	"net/http"
	"subscriptions-api-postgres/internal/handlers"
)

func New(subscriptionHandler *handlers.SubscriptionHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.HealthCheck)
	mux.HandleFunc("POST /subscriptions", subscriptionHandler.CreateSubscription)
	mux.HandleFunc("GET /subscriptions", subscriptionHandler.GetSubscriptions)

	return mux
}
