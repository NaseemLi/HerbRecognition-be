package routes

import (
	"herb-recognition-be/internal/handler/app"
	"herb-recognition-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerRecognizeRoutes(r *gin.Engine) {
	recognizeHandler := app.NewRecognizeHandler()

	recognize := r.Group("/api/recognize")
	recognize.Use(middleware.JWTAuth())
	{
		recognize.POST("/upload", recognizeHandler.UploadAndRecognize)
		recognize.GET("/history", recognizeHandler.GetHistory)
		recognize.DELETE("/history", recognizeHandler.DeleteHistory)
	}
}
