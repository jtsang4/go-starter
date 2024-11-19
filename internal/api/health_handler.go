package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jtsang4/go-stater/pkg/response"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health 检查服务健康状态
func (h *HealthHandler) Health(c *gin.Context) {
	health := map[string]interface{}{
		"status": "up",
		"db":     "up",
	}

	// 检查数据库连接
	if db, err := h.db.DB(); err != nil || db.Ping() != nil {
		health["db"] = "down"
		health["status"] = "degraded"
	}

	response.Success(c, health)
}
