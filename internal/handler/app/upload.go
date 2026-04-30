package app

import (
	"herb-recognition-be/internal/middleware"
	"herb-recognition-be/pkg/response"
	"herb-recognition-be/pkg/upload"

	"github.com/gin-gonic/gin"
)

// UploadHandler 通用上传处理器
type UploadHandler struct{}

// NewUploadHandler 创建上传处理器
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadImage 上传通用图片
func (h *UploadHandler) UploadImage(c *gin.Context) {
	_, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		response.Error(c, 401, "未登录", nil)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		response.Error(c, 400, "请上传图片文件", nil)
		return
	}

	cfg := upload.DefaultImageConfig
	cfg.UploadDir = "./uploads/common"
	cfg.URLPrefix = "/uploads/common/"

	imageURL, err := upload.UploadFile(file, cfg)
	if err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "上传成功", gin.H{"image_url": imageURL})
}
