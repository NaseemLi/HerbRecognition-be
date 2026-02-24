package routes

import (
	"herb-recognition-be/internal/handler/admin"
	"herb-recognition-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerAdminUserRoutes(r *gin.Engine) {
	adminUserHandler := admin.NewAdminUserHandler()

	admin := r.Group("/api/admin/user")
	admin.Use(middleware.JWTAuth())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("", adminUserHandler.GetUserList)          // 用户列表
		admin.POST("/role", adminUserHandler.UpdateUserRole) // 修改用户角色
	}
}
