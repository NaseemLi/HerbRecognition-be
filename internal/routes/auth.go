package routes

import (
	"herb-recognition-be/internal/handler/app"
	"herb-recognition-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(r *gin.Engine) {
	authHandler := app.NewAuthHandler()

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)

		// 需要认证的接口
		auth.Use(middleware.JWTAuth())
		{
			auth.POST("/change-password", authHandler.ChangePassword)
		}
	}
}
