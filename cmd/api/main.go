package main

import (
	"fmt"
	"net/http"
	"subscriptions-api-postgres/internal/config"
	"subscriptions-api-postgres/internal/database"
	"subscriptions-api-postgres/internal/handlers"
	"subscriptions-api-postgres/internal/repository"
	"subscriptions-api-postgres/internal/router"
	"subscriptions-api-postgres/internal/service"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		fmt.Println("db connection error:", err)
		return
	}
	defer db.Close()
	fmt.Println("db connection successfully")

	subscriptionRepository := repository.NewSubscriptionsRepository(db)
	subscriptionService := service.NewSubscriptionsService(subscriptionRepository)
	subscriptionHandler := handlers.NewSubscriptionsHandler(subscriptionService)

	appRouter := router.New(subscriptionHandler)

	address := ":" + cfg.Server.Port

	fmt.Println("server is running on", address)
	err = http.ListenAndServe(address, appRouter)
	if err != nil {
		fmt.Println("server error:", err)
		return
	}

}
