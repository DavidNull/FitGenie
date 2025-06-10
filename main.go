package main

import (
	"log"

	"fitgenie/internal/api"
	"fitgenie/internal/config"
	"fitgenie/internal/database"
	"fitgenie/internal/services"
)

func main() {
	cfg := config.Load()

	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	colorService := services.NewColorTheoryService()
	styleService := services.NewStyleService()
	aiService := services.NewAIService(colorService, styleService)

	router := api.SetupRoutes(colorService, styleService, aiService)

	log.Printf("FitGenie server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
