package routes

import (
	"herb-recognition-be/internal/handler/app"
	"herb-recognition-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerTicketRoutes(r *gin.Engine) {
	ticketHandler := app.NewTicketHandler()

	ticket := r.Group("/api/ticket")
	ticket.Use(middleware.JWTAuth())
	{
		ticket.POST("", ticketHandler.CreateTicket)       // 提交工单
		ticket.GET("", ticketHandler.GetMyTickets)          // 我的工单列表
		ticket.GET("/:id", ticketHandler.GetTicketDetail)   // 工单详情
	}
}
