package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.HEAD("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// 认证路由
	registerAuthRoutes(r)

	// 用户路由
	registerUserRoutes(r)

	// 识别路由
	registerRecognizeRoutes(r)

	// 搜索路由
	registerSearchRoutes(r)

	// 管理端药材路由
	registerAdminHerbRoutes(r)

	// 管理端用户路由
	registerAdminUserRoutes(r)

	// 工单路由
	registerTicketRoutes(r)

	// 管理端工单调路由
	registerAdminTicketRoutes(r)

	// 通用上传路由
	registerUploadRoutes(r)
}
