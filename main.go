package main

import (
	"log"

	"fitgenie/internal/api" //  handles HTTP routing and endpoints
	"fitgenie/internal/config" // application settings
	"fitgenie/internal/database" //  manages DB connection and migrations
	"fitgenie/internal/services" // logic for color theory, etc.
)

func main() {
	cfg := config.Load()

	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	colorService := services.NewColorTheoryService()
	router := api.SetupColorTheoryRoutes(colorService)

	log.Printf("FitGenie Color Theory server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
