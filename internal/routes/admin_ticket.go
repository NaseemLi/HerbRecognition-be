package routes

import (
	"herb-recognition-be/internal/handler/admin"
	"herb-recognition-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerAdminTicketRoutes(r *gin.Engine) {
	adminTicketHandler := admin.NewAdminTicketHandler()

	admin := r.Group("/api/admin/ticket")
	admin.Use(middleware.JWTAuth())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("", adminTicketHandler.GetTicketList)                   // 工单列表
		admin.GET("/:id", adminTicketHandler.GetTicketDetail)             // 工单详情
		admin.POST("/:id/reply", adminTicketHandler.ReplyTicket)          // 回复工单
		admin.POST("/:id/status", adminTicketHandler.UpdateTicketStatus)  // 修改状态
	}
}
