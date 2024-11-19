package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jtsang4/go-stater/config"
	"github.com/jtsang4/go-stater/pkg/auth"
	"github.com/jtsang4/go-stater/pkg/logger"
	"go.uber.org/zap"
)

func AuthMiddleware(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Logger.Info("missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			logger.Logger.Info("invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(bearerToken[1], cfg)
		if err != nil {
			logger.Logger.Error("failed to parse token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store user information from claims in context
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
