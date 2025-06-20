package api

import (
	"fitgenie/internal/api/handlers"
	"fitgenie/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	colorService *services.ColorTheoryService,
	styleService *services.StyleService,
	aiService *services.AIService,
) *gin.Engine {

	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	userHandler := handlers.NewUserHandler()
	outfitHandler := handlers.NewOutfitHandler(aiService)

	// API version 1
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "healthy",
				"service": "FitGenie API",
				"version": "1.0.0",
			})
		})

		// User routes
		users := v1.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/:userId", userHandler.GetUser)
			users.PUT("/:userId", userHandler.UpdateUser)
			users.DELETE("/:userId", userHandler.DeleteUser)

			// User profile 
			users.POST("/:userId/style-profile", userHandler.CreateStyleProfile)
			users.GET("/:userId/style-profile", userHandler.GetStyleProfile)
			users.POST("/:userId/color-profile", userHandler.CreateColorProfile)
			users.GET("/:userId/color-profile", userHandler.GetColorProfile)

			// Favorite outfits 
			users.POST("/:userId/favorites/:outfitId", userHandler.AddFavoriteOutfit)
			users.DELETE("/:userId/favorites/:outfitId", userHandler.RemoveFavoriteOutfit)
			users.GET("/:userId/favorites", userHandler.ListFavoriteOutfits)
		}

		// Clothing items 
		clothing := v1.Group("/users/:userId/clothing")
		{
			clothing.POST("", outfitHandler.CreateClothingItem)
			clothing.GET("", outfitHandler.ListClothingItems)
			clothing.GET("/:id", outfitHandler.GetClothingItem)
			clothing.PUT("/:id", outfitHandler.UpdateClothingItem)
			clothing.DELETE("/:id", outfitHandler.DeleteClothingItem)
		}

		// Outfits routes
		outfits := v1.Group("/users/:userId/outfits")
		{
			outfits.POST("", outfitHandler.CreateOutfit)
			outfits.GET("", outfitHandler.ListOutfits)
			outfits.GET("/:id", outfitHandler.GetOutfit)
			outfits.DELETE("/:id", outfitHandler.DeleteOutfit)
		}

		// AI recommendations routes
		recommendations := v1.Group("/users/:userId/recommendations")
		{
			recommendations.POST("/outfits", outfitHandler.GetOutfitRecommendations)
			recommendations.GET("/personalized", outfitHandler.GetPersonalizedRecommendations)
		}

		// Recommendation interaction routes
		recInteractions := v1.Group("/recommendations")
		{
			recInteractions.PUT("/:id/viewed", outfitHandler.MarkRecommendationViewed)
			recInteractions.PUT("/:id/accepted", outfitHandler.AcceptRecommendation)
		}

		// Color theory routes
		colorTheory := v1.Group("/color-theory")
		{
			colorTheory.GET("/seasons", func(c *gin.Context) {
				seasons := []string{"spring", "summer", "autumn", "winter"}
				c.JSON(200, gin.H{"seasons": seasons})
			})

			colorTheory.GET("/harmonies", func(c *gin.Context) {
				harmonies := []string{
					"monochromatic", "analogous", "complementary",
					"triadic", "split-complementary", "tetradic",
				}
				c.JSON(200, gin.H{"harmonies": harmonies})
			})
		}

		// Style analysis routes
		styleAnalysis := v1.Group("/style-analysis")
		{
			styleAnalysis.GET("/categories", func(c *gin.Context) {
				categories := []string{
					"casual", "business", "formal", "bohemian", "minimalist",
					"vintage", "sporty", "romantic", "edgy", "classic",
				}
				c.JSON(200, gin.H{"categories": categories})
			})

			styleAnalysis.GET("/body-types", func(c *gin.Context) {
				bodyTypes := []string{
					"pear", "apple", "hourglass", "rectangle", "inverted-triangle",
				}
				c.JSON(200, gin.H{"body_types": bodyTypes})
			})

			styleAnalysis.GET("/occasions", func(c *gin.Context) {
				occasions := []string{
					"work", "casual", "formal", "party", "date", "travel",
					"workout", "beach", "wedding", "interview",
				}
				c.JSON(200, gin.H{"occasions": occasions})
			})
		}

		// Analytics routes (for future implementation)
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/user/:userId/wardrobe-stats", func(c *gin.Context) {
				// Placeholder for wardrobe analytics
				c.JSON(200, gin.H{
					"message": "Wardrobe analytics endpoint - coming soon",
				})
			})

			analytics.GET("/user/:userId/style-insights", func(c *gin.Context) {
				// Placeholder for style insights
				c.JSON(200, gin.H{
					"message": "Style insights endpoint - coming soon",
				})
			})
		}
	}

	return router
}
