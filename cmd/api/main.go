package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"

	"github.com/dokkiitech/grumble-back/internal/api"
	"github.com/dokkiitech/grumble-back/internal/config"
	"github.com/dokkiitech/grumble-back/internal/controller"
	"github.com/dokkiitech/grumble-back/internal/controller/middleware"
	sharedservice "github.com/dokkiitech/grumble-back/internal/domain/shared/service"
	"github.com/dokkiitech/grumble-back/internal/infrastructure"
	"github.com/dokkiitech/grumble-back/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/api/option"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create logger
	logger := infrastructure.NewLogger()
	logger.Info("Starting Grumble API server")

	// Connect to database
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Test database connection
	if err := dbPool.Ping(ctx); err != nil {
		logger.Error("Failed to ping database", "error", err)
		log.Fatalf("Failed to ping database: %v", err)
	}
	logger.Info("Database connection established")

	// Migrations must be applied separately via `go run ./cmd/migrate` before starting the API

	// Initialize Firebase authentication
	var firebaseOpts []option.ClientOption
	if cfg.FirebaseCredentialsFile != "" {
		firebaseOpts = append(firebaseOpts, option.WithCredentialsFile(cfg.FirebaseCredentialsFile))
	}

	var firebaseCfg *firebase.Config
	if cfg.FirebaseProjectID != "" {
		firebaseCfg = &firebase.Config{ProjectID: cfg.FirebaseProjectID}
	}

	firebaseApp, err := firebase.NewApp(ctx, firebaseCfg, firebaseOpts...)
	if err != nil {
		logger.Error("Failed to initialize Firebase app", "error", err)
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		logger.Error("Failed to create Firebase auth client", "error", err)
		log.Fatalf("Failed to create Firebase auth client: %v", err)
	}
	logger.Info("Firebase authentication client initialized")

	// Initialize repositories
	grumbleRepo := infrastructure.NewPostgresGrumbleRepository(dbPool)
	userRepo := infrastructure.NewPostgresUserRepository(dbPool)
	vibeRepo := infrastructure.NewPostgresVibeRepository(dbPool)

	// Initialize use cases
	grumblePostUC := usecase.NewGrumblePostUseCase(grumbleRepo)
	timelineGetUC := usecase.NewTimelineGetUseCase(grumbleRepo)
	authAnonymousUC := usecase.NewAuthAnonymousUseCase(userRepo)
	userQueryUC := usecase.NewUserQueryUseCase(userRepo)
	purifyService := sharedservice.NewPurifyService(cfg.PurificationThreshold)
	virtueService := sharedservice.NewVirtueService()
	vibeAddUC := usecase.NewVibeAddUseCase(grumbleRepo, vibeRepo, userRepo, purifyService, virtueService)

	// Initialize presenters
	grumblePresenter := controller.NewGrumblePresenter()
	timelinePresenter := controller.NewTimelinePresenter(grumblePresenter)

	// Initialize controllers
	grumbleController := controller.NewGrumbleController(grumblePostUC, grumblePresenter, logger)
	timelineController := controller.NewTimelineController(timelineGetUC, timelinePresenter, logger)
	authController := controller.NewAuthController(authAnonymousUC, userQueryUC, logger)
	vibeController := controller.NewVibeController(vibeAddUC, logger)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authClient, authAnonymousUC, logger)

	// Create server implementation that combines all controllers
	serverImpl := api.NewServerImpl(grumbleController, timelineController, authController, vibeController)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Global middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		logger.Info("Request completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
		)
	})

	// Apply auth middleware to protected routes
	// For MVP, we'll apply it globally for simplicity
	router.Use(authMiddleware.Authenticate())

	// Register OpenAPI routes with /api/v1 prefix
	api.RegisterHandlersWithOptions(router, serverImpl, api.GinServerOptions{
		BaseURL: "/api/v1",
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		logger.Info("HTTP server listening", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}
