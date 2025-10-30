package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration.
type Config struct {
	// Server
	HTTPAddr string

	// Database
	DatabaseURL string

	// Business Rules
	PurificationThreshold int

	// Performance
	DBMaxConns int
	DBMinConns int
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		HTTPAddr:              getEnv("GRUMBLE_HTTP_ADDR", ":8080"),
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		PurificationThreshold: getEnvInt("PURIFICATION_THRESHOLD", 10),
		DBMaxConns:            getEnvInt("DB_MAX_CONNS", 25),
		DBMinConns:            getEnvInt("DB_MIN_CONNS", 5),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
