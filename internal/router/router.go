package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jtsang4/go-stater/config"
	"github.com/jtsang4/go-stater/internal/api"
	"github.com/jtsang4/go-stater/internal/middleware"
)

func SetupRouter(r *gin.Engine, userHandler *api.UserHandler, cfg *config.Config, healthHandler *api.HealthHandler) {
	// Health check route
	r.GET("/health", healthHandler.Health)

	// Public routes
	public := r.Group("/api/v1")
	{
		public.POST("/users/register", userHandler.Register)
		public.POST("/users/login", userHandler.Login)
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWT))
	{
		protected.GET("/users/:id", userHandler.GetUser)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)
	}
}
