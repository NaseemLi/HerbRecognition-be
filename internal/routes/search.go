package routes

import (
	"herb-recognition-be/internal/handler/app"
	"herb-recognition-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerSearchRoutes(r *gin.Engine) {
	searchHandler := app.NewSearchHandler()

	search := r.Group("/api/herb")
	search.Use(middleware.JWTAuth())
	{
		search.GET("", searchHandler.GetAll)        // 获取药材列表
		search.GET("/search", searchHandler.Search) // 搜索药材
		search.GET("/:id", searchHandler.GetDetail) // 药材详情
	}
}
