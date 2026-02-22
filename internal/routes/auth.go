package routes

import (
	"herb-recognition-be/internal/handler/app"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(r *gin.Engine) {
	authHandler := app.NewAuthHandler()

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}
