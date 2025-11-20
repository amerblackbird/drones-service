package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "drones/configs"
	httpHandler "drones/internal/adapters/http"
	natsadapter "drones/internal/adapters/nats"

	"github.com/gorilla/mux"

	"drones/internal/adapters/logger"
	"drones/internal/adapters/postgres"
	"drones/internal/adapters/redis"
	"drones/internal/core/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger, err := logger.NewProductionZapLogger()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// Initialize database connection
	db, err := postgres.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Close db connection on exit
	defer func() {
		if err := db.Close(); err != nil {
			appLogger.Error("Failed to close database connection", "error", err)
		} else {
			appLogger.Info("Database connection closed")
		}
	}()

	// Cache service initialization
	cacheClient := redis.NewRedisClient(cfg.Redis)

	result, err := cacheClient.Ping(context.Background()).Result() // Test connection
	if err != nil {
		appLogger.Error("Failed to ping Redis", "error", err)
	} else {
		appLogger.Info("Redis ping successful", "result", result)
	}

	cacheService := redis.NewRedisCacheService(cacheClient, appLogger)

	defer func() {
		if err := cacheService.Close(); err != nil {
			appLogger.Error("Failed to close Redis client", "error", err)
		} else {
			appLogger.Info("Redis client closed")
		}
	}()

	// Initialize repositories
	usersRepo := postgres.NewUserRepository(db, appLogger)
	// loginRepo := postgres.NewLoginsRepository(db, appLogger)
	dronesRepo := postgres.NewDronesRepository(db, appLogger)
	ordersRepo := postgres.NewOrdersRepository(db, appLogger)
	// activityLogsRepo := postgres.NewActivityLogsRepository(db, appLogger)
	// auditLogsRepo := postgres.NewAuditLogsRepository(db, appLogger)

	natsEventPublisher := natsadapter.NewEventPublisher(cfg.NATS, appLogger)
	natsEventConsumer := natsadapter.NewEventConsumer(cfg.NATS, appLogger)

	appLogger.Info("NATS event publisher started successfully")

	// Initialize services
	usersService := services.NewUserRepository(usersRepo, natsEventPublisher, cacheService, appLogger)
	dronesService := services.NewDronesService(dronesRepo, cacheService, natsEventPublisher, appLogger)

	ordersService := services.NewOrdersService(ordersRepo, dronesService, cacheService, natsEventPublisher, appLogger)
	tokenService := services.NewJWTService(&cfg.Jwt)
	authService := services.NewAuthService(usersService, tokenService, cfg.Jwt, appLogger)
	// activityLogsService := services.NewActivityLogsService(activityLogsRepo, cacheService, natsEventPublisher, appLogger)
	// auditLogsService := services.NewAuditLogsService(auditLogsRepo, cacheService, natsEventPublisher, appLogger)

	natsEventHandlers := natsadapter.NewEventHandlers(dronesService, appLogger)
	natsEventHandlers.RegisterHandlers(natsEventConsumer)

	// Initialize HTTP handler
	httpHandlerInstance := httpHandler.NewHTTPHandler(authService, ordersService, dronesService, natsEventPublisher, appLogger, cfg.Server.ApiPrefix)

	// Setup routes
	r := mux.NewRouter()

	httpHandlerInstance.SetupRoutes(r)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Start event consumers
	ctx := context.Background()

	// Start NATS event consumer
	go func() {
		// Start NATS event consumer
		if err := natsEventConsumer.Start(ctx); err != nil {
			log.Fatalf("HTTP Server starting: %v", err)
		}

		// Start NATS event publisher
		if err := natsEventPublisher.Start(); err != nil {
			log.Fatalf("Failed to start NATS event publisher: %v", err)
		}
	}()

	// Start server in a goroutine
	go func() {
		appLogger.Info("Server starting", "address", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Server shutting down...")

	// Create a deadline to wait for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop event consumers and publishers
	if err := natsEventConsumer.Stop(); err != nil {
		appLogger.Error("Error stopping NATS event consumer", "error", err)
	}

	if err := natsEventPublisher.Stop(); err != nil {
		appLogger.Error("Error stopping NATS event publisher", "error", err)
	}

	// Close cache service
	if err := cacheService.Close(); err != nil {
		appLogger.Error("Error closing cache service", "error", err)
	}

	// Close repos
	if err := usersRepo.Close(); err != nil {
		appLogger.Error("Error closing usersRepo", "error", err)
	}
	if err := dronesRepo.Close(); err != nil {
		appLogger.Error("Error closing dronesRepo", "error", err)
	}
	if err := ordersRepo.Close(); err != nil {
		appLogger.Error("Error closing ordersRepo", "error", err)
	}

	// Close database connection
	if err := db.Close(); err != nil {
		appLogger.Error("Error closing database connection", "error", err)
	}

	// Shutdown HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	appLogger.Info("Server exited")

}
