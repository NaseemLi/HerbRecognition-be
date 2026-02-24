package routes

import (
	"herb-recognition-be/internal/handler/app"
	"herb-recognition-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerUserRoutes(r *gin.Engine) {
	userHandler := app.NewUserHandler()

	user := r.Group("/api/user")
	user.Use(middleware.JWTAuth())
	{
		user.GET("/profile", userHandler.GetProfile)    // 获取用户资料
		user.PUT("/profile", userHandler.UpdateProfile) // 更新用户资料
		user.POST("/avatar", userHandler.UploadAvatar)  // 上传头像
	}
}
