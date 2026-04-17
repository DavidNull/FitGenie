package api

import (
	"fitgenie/internal/api/handlers"
	"fitgenie/internal/config"
	"fitgenie/internal/repository"
	"fitgenie/internal/services"
	"fitgenie/pkg/database"
	"fitgenie/pkg/logger"
	"fitgenie/pkg/middleware"
	"fitgenie/pkg/storage"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *database.Connection, log *logger.Logger, cfg *config.Config) *gin.Engine {
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

	var s3Client *storage.S3Client
	if cfg.S3Endpoint != "" {
		var err error
		s3Client, err = storage.NewS3Client(storage.S3Config{
			Endpoint:        cfg.S3Endpoint,
			Region:          cfg.S3Region,
			Bucket:          cfg.S3Bucket,
			AccessKeyID:     cfg.S3AccessKeyID,
			SecretAccessKey: cfg.S3SecretAccessKey,
			UsePathStyle:    cfg.S3UsePathStyle,
		})
		if err != nil {
			log.Error("failed to initialize S3 client", "error", err)
		}
	}

	uploadHandler := handlers.NewUploadHandler(s3Client, log)

	deviceAuth := middleware.DeviceAuthMiddleware(userRepo, log)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/upload", deviceAuth, uploadHandler.UploadImage)
		v1.GET("/images/:path", deviceAuth, uploadHandler.GetImageURL)

		users := v1.Group("/users")
		users.Use(deviceAuth)
		{
			users.POST("", userHandler.CreateUser)
			users.GET("/me", userHandler.GetCurrentUser)
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
		clothing.Use(deviceAuth)
		{
			clothing.GET("", clothingHandler.ListClothing)
			clothing.POST("", clothingHandler.CreateClothing)
			clothing.GET("/:id", clothingHandler.GetClothing)
			clothing.PUT("/:id", clothingHandler.UpdateClothing)
			clothing.DELETE("/:id", clothingHandler.DeleteClothing)
		}

		outfits := v1.Group("/outfits")
		outfits.Use(deviceAuth)
		{
			outfits.POST("", outfitHandler.CreateOutfit)
			outfits.GET("/:id", outfitHandler.GetOutfit)
			outfits.PUT("/:id", outfitHandler.UpdateOutfit)
			outfits.DELETE("/:id", outfitHandler.DeleteOutfit)
		}

		userOutfits := v1.Group("/users/:userId/outfits")
		userOutfits.Use(deviceAuth)
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
