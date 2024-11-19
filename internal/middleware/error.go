package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jtsang4/go-stater/pkg/logger"
	"go.uber.org/zap"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 实现 error 接口
func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last()
		logger.Logger.Error("request error", zap.Error(err))

		if appErr, ok := err.Err.(*AppError); ok {
			c.JSON(appErr.Code, appErr)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
	}
}
