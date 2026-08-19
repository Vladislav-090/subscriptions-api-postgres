package main

import (
	"log/slog"
	"net/http"
	"os"
	"subscriptions-api-postgres/internal/auth"
	"subscriptions-api-postgres/internal/config"
	"subscriptions-api-postgres/internal/database"
	"subscriptions-api-postgres/internal/handlers"
	"subscriptions-api-postgres/internal/repository"
	"subscriptions-api-postgres/internal/router"
	"subscriptions-api-postgres/internal/service"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return
	}
	defer db.Close()
	slog.Info("database connected")

	subscriptionRepository := repository.NewSubscriptionsRepository(db)
	subscriptionService := service.NewSubscriptionsService(subscriptionRepository)
	subscriptionHandler := handlers.NewSubscriptionsHandler(subscriptionService)

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.TTL)

	usersRepository := repository.NewUsersRepository(db)
	revokedTokensRepository := repository.NewRevokedTokensRepository(db)
	authService := service.NewAuthService(usersRepository, revokedTokensRepository, jwtManager)
	authHandler := handlers.NewAuthHandler(authService)

	appRouter := router.New(subscriptionHandler, authHandler, jwtManager, authService)

	address := ":" + cfg.Server.Port

	slog.Info("server starting", "address", address)
	err = http.ListenAndServe(address, appRouter)
	if err != nil {
		slog.Error("server failed", "error", err)
		return
	}

}
