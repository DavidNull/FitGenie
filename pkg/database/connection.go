package database

import (
	"fmt"
	"time"

	"fitgenie/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Connection struct {
	*gorm.DB
}

func NewConnection(databaseURL string) (*Connection, error) {
	gormLogger := logger.Default.LogMode(logger.Info)

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return &Connection{DB: db}, nil
}

func (c *Connection) Migrate() error {
	models := []interface{}{
		&models.User{},
		&models.StyleProfile{},
		&models.ColorProfile{},
		&models.ClothingItem{},
		&models.Outfit{},
		&models.OutfitRecommendation{},
		&models.FavoriteOutfit{},
	}

	for _, model := range models {
		if err := c.DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	return nil
}

func (c *Connection) Health() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (c *Connection) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
