package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
	ColorAnalysisEnabled bool
	StyleAnalysisEnabled bool
	MaxImageSize      int64
	AllowedImageTypes []string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://fitgenie:fitgenie123@localhost:5432/fitgenie?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "default-secret-key"),
		ColorAnalysisEnabled: getBoolEnv("COLOR_ANALYSIS_ENABLED", true),
		StyleAnalysisEnabled: getBoolEnv("STYLE_ANALYSIS_ENABLED", true),
		MaxImageSize:      getInt64Env("MAX_IMAGE_SIZE", 5242880), // 5MB
		AllowedImageTypes: []string{"jpg", "jpeg", "png", "webp"},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return fallback
}

func getInt64Env(key string, fallback int64) int64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return fallback
}
