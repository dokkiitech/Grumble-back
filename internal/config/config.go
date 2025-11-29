package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Config holds all application configuration.
type Config struct {
	// Server
	HTTPAddr string

	// Database
	DatabaseURL string

	// Authentication
	FirebaseProjectID       string
	FirebaseCredentialsFile string

	// Business Rules
	PurificationThresholdDefault   int
	PurificationThresholdMin       int
	PurificationThresholdMax       int
	BodhisattvaRankingLimitDefault int
	BodhisattvaRankingLimitMin     int
	BodhisattvaRankingLimitMax     int

	// HTTP
	CORSAllowedOrigins []string
	GinMode            string

	// Performance
	DBMaxConns int
	DBMinConns int

	// Content Moderation
	GeminiAPIKey string
	GeminiModel  string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		HTTPAddr:                       getEnv("GRUMBLE_HTTP_ADDR", ":8080"),
		DatabaseURL:                    os.Getenv("DATABASE_URL"),
		FirebaseProjectID:              os.Getenv("FIREBASE_PROJECT_ID"),
		FirebaseCredentialsFile:        os.Getenv("FIREBASE_CREDENTIALS_FILE"),
		CORSAllowedOrigins:             getEnvStringSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:8081", "http://localhost:19006"}),
		GinMode:                        getEnv("GIN_MODE", gin.ReleaseMode),
		PurificationThresholdDefault:   getEnvInt("PURIFICATION_THRESHOLD_DEFAULT", 10),
		PurificationThresholdMin:       getEnvInt("PURIFICATION_THRESHOLD_MIN", 1),
		PurificationThresholdMax:       getEnvInt("PURIFICATION_THRESHOLD_MAX", 1000),
		BodhisattvaRankingLimitDefault: getEnvInt("BODHISATTVA_RANKING_LIMIT_DEFAULT", 10),
		BodhisattvaRankingLimitMin:     getEnvInt("BODHISATTVA_RANKING_LIMIT_MIN", 1),
		BodhisattvaRankingLimitMax:     getEnvInt("BODHISATTVA_RANKING_LIMIT_MAX", 100),
		DBMaxConns:                     getEnvInt("DB_MAX_CONNS", 25),
		DBMinConns:                     getEnvInt("DB_MIN_CONNS", 5),
		GeminiAPIKey:                   os.Getenv("GEMINI_API_KEY"),
		GeminiModel:                    getEnv("GEMINI_MODEL", "gemini-2.5-flash-lite"),
	}

	if cfg.FirebaseCredentialsFile == "" {
		const defaultSecretsFile = "firebase_secrets.json"
		if _, err := os.Stat(defaultSecretsFile); err == nil {
			if absPath, err := filepath.Abs(defaultSecretsFile); err == nil {
				cfg.FirebaseCredentialsFile = absPath
			} else {
				cfg.FirebaseCredentialsFile = defaultSecretsFile
			}
		}
	}

	if cfg.FirebaseCredentialsFile != "" {
		if projectID := extractProjectID(cfg.FirebaseCredentialsFile); projectID != "" {
			cfg.FirebaseProjectID = projectID
		}
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
	if intVal, err := strconv.Atoi(getEnv(key, "")); err == nil {
		return intVal
	}
	return defaultValue
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		var result []string
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

func extractProjectID(credentialsPath string) string {
	data, err := os.ReadFile(credentialsPath)
	if err != nil {
		return ""
	}
	var cred struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal(data, &cred); err != nil {
		return ""
	}
	return strings.TrimSpace(cred.ProjectID)
}
