package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dokkiitech/grumble-back/internal/config"
	sharedservice "github.com/dokkiitech/grumble-back/internal/domain/shared/service"
	"github.com/dokkiitech/grumble-back/internal/infrastructure"
	"github.com/dokkiitech/grumble-back/internal/job"
	"github.com/dokkiitech/grumble-back/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	mode := flag.String("mode", "cron", "batch mode: cron|purge-expired")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := infrastructure.NewLogger()
	logger.Info("Starting batch", "mode", *mode)

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

	// Initialize domain service
	eventTimeService := sharedservice.NewEventTimeService()

	grumbleRepo := infrastructure.NewPostgresGrumbleRepository(dbPool, eventTimeService)
	purgeUC := usecase.NewPurgeExpiredUseCase(grumbleRepo, logger)

	switch *mode {
	case "purge-expired":
		// Run once
		if _, err := purgeUC.Purge(ctx); err != nil {
			logger.Error("purge-expired failed", "error", err)
			os.Exit(1)
		}
		logger.Info("purge-expired completed")
	case "cron":
		// Start scheduler
		scheduler := job.NewCronScheduler(job.NewPurgeExpiredJob(purgeUC, logger), logger)
		if err := scheduler.Start(); err != nil {
			logger.Error("Scheduler start failed", "error", err)
			os.Exit(1)
		}
		defer scheduler.Stop()

		// Wait for termination
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		logger.Info("Shutting down batch")
	default:
		logger.Error("Unknown mode", "mode", *mode)
		os.Exit(2)
	}
}
