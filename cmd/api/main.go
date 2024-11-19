package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jtsang4/go-stater/config"
	"github.com/jtsang4/go-stater/internal/api"
	"github.com/jtsang4/go-stater/internal/middleware"
	"github.com/jtsang4/go-stater/internal/model"
	"github.com/jtsang4/go-stater/internal/repository"
	"github.com/jtsang4/go-stater/internal/router"
	"github.com/jtsang4/go-stater/internal/service"
	"github.com/jtsang4/go-stater/pkg/cache"
	"github.com/jtsang4/go-stater/pkg/database"
	"github.com/jtsang4/go-stater/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger.InitLogger(cfg.Logger)
	defer logger.Logger.Sync()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database
	db := database.InitDB(cfg.Database)

	// Initialize Redis
	redisCache := cache.NewRedisCache(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	// Auto migrate models
	if err := db.AutoMigrate(&model.User{}); err != nil {
		logger.Logger.Fatal("Failed to auto migrate database", zap.Error(err))
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, redisCache)

	// Initialize handlers
	userHandler := api.NewUserHandler(userService, cfg)
	healthHandler := api.NewHealthHandler(db)

	// Create Gin engine
	r := gin.Default()

	// Setup middleware
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.CORSMiddleware())

	// Setup routes
	router.SetupRouter(r, userHandler, cfg, healthHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		logger.Logger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	// Close database connection
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}

	logger.Logger.Info("Server exiting")
}
