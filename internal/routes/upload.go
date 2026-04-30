package routes

import (
	"herb-recognition-be/internal/handler/app"
	"herb-recognition-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerUploadRoutes(r *gin.Engine) {
	uploadHandler := app.NewUploadHandler()

	upload := r.Group("/api/upload")
	upload.Use(middleware.JWTAuth())
	{
		upload.POST("/image", uploadHandler.UploadImage)
	}
}
