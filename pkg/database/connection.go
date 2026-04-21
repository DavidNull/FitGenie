package database

import (
	"embed"
	"fitgenie/internal/models"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Connection struct {
	*gorm.DB
	dbURL string
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

	return &Connection{DB: db, dbURL: databaseURL}, nil
}

func (c *Connection) Migrate() error {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}
	defer d.Close()

	m, err := migrate.NewWithSourceInstance("iofs", d, c.dbURL)
	if err != nil {
		// Fallback: use GORM AutoMigrate if golang-migrate fails
		return c.autoMigrate()
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (c *Connection) autoMigrate() error {
	return c.DB.AutoMigrate(
		&models.User{},
		&models.StyleProfile{},
		&models.ColorProfile{},
		&models.ClothingItem{},
		&models.Outfit{},
		&models.OutfitRecommendation{},
		&models.FavoriteOutfit{},
	)
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
