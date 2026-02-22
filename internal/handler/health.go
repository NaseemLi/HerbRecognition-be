package handler

import (
	"herb-recognition-be/internal/repository"
	"herb-recognition-be/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check 健康检查接口
func (h *HealthHandler) Check(c *gin.Context) {
	// 检查数据库连接
	sqlDB, err := repository.DB.DB()
	if err != nil {
		logger.Errorf("数据库连接检查失败：%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "数据库连接失败",
			"data":    nil,
		})
		return
	}

	// Ping 数据库
	if err := sqlDB.Ping(); err != nil {
		logger.Errorf("数据库 Ping 失败：%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "数据库不可用",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "服务正常",
		"data": gin.H{
			"status": "healthy",
			"db":     "connected",
		},
	})
}
