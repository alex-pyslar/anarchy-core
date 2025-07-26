package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"anarchy-core/internal/api"
	"anarchy-core/internal/api/handler"
	"anarchy-core/internal/auth"
	"anarchy-core/internal/config"
	"anarchy-core/internal/database"
	"anarchy-core/internal/repository/postgres"
	"anarchy-core/internal/service"
	"anarchy-core/internal/util"

	"github.com/labstack/echo/v4"
)

func main() {
	// 1. Initialize Logger
	logger := util.NewLogger()

	// 2. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	// 3. Initialize Database Connection
	db, err := database.InitPostgresDB(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("Failed to initialize database: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database connection: %v", err)
		}
		logger.Info("Database connection closed.")
	}()

	// 4. Initialize Repositories
	userRepo := postgres.NewUserRepositoryPostgres(db)
	playerMovementRepo := postgres.NewPlayerMovementRepositoryPostgres(db)

	// 5. Initialize JWT Manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecretKey)

	// 6. Initialize Services
	authService := service.NewAuthService(userRepo, jwtManager, logger)
	playerService := service.NewPlayerService(playerMovementRepo, logger)
	websocketService := service.NewWebSocketService(logger)

	// Start WebSocket service in a goroutine
	go websocketService.Run()

	// 7. Initialize Handlers
	authHandler := handler.NewAuthHandler(authService, logger)
	playerMovementHandler := handler.NewPlayerMovementHandler(playerService, websocketService, jwtManager, logger)

	// 8. Initialize Echo Web Server
	e := echo.New()

	// 9. Setup Routes
	api.SetupRouter(e, authHandler, playerMovementHandler, jwtManager, logger)

	// 10. Start Server in a goroutine
	go func() {
		logger.Info("Starting server on port %s", cfg.AppPort)
		if err := e.Start(":" + cfg.AppPort); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start: %v", err)
			os.Exit(1)
		}
	}()

	// 11. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	// Listen for Ctrl+C (SIGINT) and graceful shutdown (SIGTERM)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until a signal is received

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
	} else {
		logger.Info("Server gracefully stopped.")
	}
}
