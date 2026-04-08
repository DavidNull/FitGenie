package api

import (
	"fitgenie/internal/api/handlers"
	"fitgenie/internal/repository"
	"fitgenie/internal/services"
	"fitgenie/pkg/database"
	"fitgenie/pkg/logger"
	"fitgenie/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *database.Connection, log *logger.Logger) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware(log))
	router.Use(middleware.PrometheusMiddleware())

	userRepo := repository.NewUserRepository(db.DB)
	clothingRepo := repository.NewClothingRepository(db.DB)
	outfitRepo := repository.NewOutfitRepository(db.DB)

	colorService := services.NewColorTheoryService()
	styleService := services.NewStyleService()
	aiService := services.NewAIService(colorService, styleService)

	userHandler := handlers.NewUserHandler(userRepo, log)
	clothingHandler := handlers.NewClothingHandler(clothingRepo, log)
	outfitHandler := handlers.NewOutfitHandler(outfitRepo, clothingRepo, userRepo, aiService, log)
	colorHandler := handlers.NewColorHandler(colorService, log)

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/:userId", userHandler.GetUser)
			users.PUT("/:userId", userHandler.UpdateUser)
			users.DELETE("/:userId", userHandler.DeleteUser)

			users.POST("/:userId/style-profile", userHandler.CreateStyleProfile)
			users.GET("/:userId/style-profile", userHandler.GetStyleProfile)
			users.POST("/:userId/color-profile", userHandler.CreateColorProfile)
			users.GET("/:userId/color-profile", userHandler.GetColorProfile)

			users.POST("/:userId/favorites/:outfitId", userHandler.AddFavoriteOutfit)
			users.DELETE("/:userId/favorites/:outfitId", userHandler.RemoveFavoriteOutfit)
			users.GET("/:userId/favorites", userHandler.ListFavoriteOutfits)
		}

		clothing := v1.Group("/clothing")
		{
			clothing.POST("", clothingHandler.CreateClothing)
			clothing.GET("/:id", clothingHandler.GetClothing)
			clothing.PUT("/:id", clothingHandler.UpdateClothing)
			clothing.DELETE("/:id", clothingHandler.DeleteClothing)
			clothing.GET("", clothingHandler.ListClothing)
		}

		outfits := v1.Group("/outfits")
		{
			outfits.POST("", outfitHandler.CreateOutfit)
			outfits.GET("/:id", outfitHandler.GetOutfit)
			outfits.DELETE("/:id", outfitHandler.DeleteOutfit)
		}

		userOutfits := v1.Group("/users/:userId/outfits")
		{
			userOutfits.GET("", outfitHandler.ListOutfits)
			userOutfits.POST("/recommendations", outfitHandler.GetOutfitRecommendations)
		}

		colorTheory := v1.Group("/color-theory")
		{
			colorTheory.GET("/seasons", colorHandler.GetSeasons)
			colorTheory.GET("/harmonies", colorHandler.GetHarmonies)
			colorTheory.POST("/analyze-harmony", colorHandler.AnalyzeHarmony)
			colorTheory.POST("/recommendations", colorHandler.GetRecommendations)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	return router
}
