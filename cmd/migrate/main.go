package main

import (
	"context"
	"log"

	"github.com/dokkiitech/grumble-back/internal/config"
	"github.com/dokkiitech/grumble-back/internal/infrastructure"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Logger
	logger := infrastructure.NewLogger()
	logger.Info("Starting migration runner")

	// Connect DB
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect database", "error", err)
		log.Fatalf("DB connect error: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		logger.Error("Failed to ping database", "error", err)
		log.Fatalf("DB ping error: %v", err)
	}

	// Run migrations
	if err := infrastructure.RunMigrations(ctx, dbPool); err != nil {
		logger.Error("Migrations failed", "error", err)
		log.Fatalf("Migrations failed: %v", err)
	}

	logger.Info("Migrations completed successfully")
}
