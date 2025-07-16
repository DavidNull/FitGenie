package database

import (
	"fmt"
	"log"

	"fitgenie/internal/config"
	"fitgenie/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize sets up the database connection and runs migrations
func Initialize(cfg *config.Config) error {
	if DB != nil {
		return nil // Prevent double initialization
	}
	var err error
	gormLogger := logger.Default.LogMode(logger.Info)
	DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	err = runMigrations()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Printf("Database initialized successfully")
	return nil
}

func runMigrations() error {
	models := []interface{}{
		&models.User{},
		&models.StyleProfile{},
		&models.ColorProfile{},
		&models.ClothingItem{},
		&models.Outfit{},
		&models.OutfitRecommendation{},
	}

	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to run auto-migrations: %w", err)
		}
	}

	log.Printf("Database migrations completed successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func Health() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
