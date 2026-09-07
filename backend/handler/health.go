package handler

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"showcase-backend/pkg/response"
)

// HealthHandler 提供 /api/health（不依赖 DB 也能响应，便于先跑通骨架）。
type HealthHandler struct {
	DB *gorm.DB
}

func (h *HealthHandler) Health(c *gin.Context) {
	dbStatus := "not_configured"
	if h.DB != nil {
		sqlDB, err := h.DB.DB()
		if err != nil {
			dbStatus = "error"
		} else if err := sqlDB.Ping(); err != nil {
			dbStatus = "error: " + err.Error()
		} else {
			dbStatus = "connected"
		}
	}
	response.OK(c, gin.H{
		"service": "showcase-backend",
		"status":  "ok",
		"db":      dbStatus,
	})
}
