package routes

import (
	"herb-recognition-be/internal/handler/admin"
	"herb-recognition-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerAdminHerbRoutes(r *gin.Engine) {
	adminHerbHandler := admin.NewAdminHerbHandler()

	admin := r.Group("/api/admin/herb")
	admin.Use(middleware.JWTAuth())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("", adminHerbHandler.GetList)                         // 药材列表
		admin.POST("", adminHerbHandler.CreateHerb)                     // 新增药材
		admin.PUT("", adminHerbHandler.UpdateHerb)                      // 编辑药材
		admin.DELETE("", adminHerbHandler.DeleteHerb)                   // 删除药材
		admin.DELETE("/batch", adminHerbHandler.BatchDelete)            // 批量删除
		admin.POST("/upload-image", adminHerbHandler.UploadAndSetImage) // 上传图片
	}
}
