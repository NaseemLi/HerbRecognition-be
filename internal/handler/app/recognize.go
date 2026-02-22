package app

import (
	"herb-recognition-be/internal/middleware"
	"herb-recognition-be/internal/service"
	"herb-recognition-be/pkg/response"

	"github.com/gin-gonic/gin"
)

// RecognizeHandler 识别处理器
type RecognizeHandler struct {
	recognizeService *service.RecognizeService
}

// NewRecognizeHandler 创建识别处理器
func NewRecognizeHandler() *RecognizeHandler {
	return &RecognizeHandler{
		recognizeService: &service.RecognizeService{},
	}
}

// UploadAndRecognize 上传并识别图片
func (h *RecognizeHandler) UploadAndRecognize(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, 401, "未登录", nil)
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		response.Error(c, 400, "请上传图片文件", nil)
		return
	}

	// 上传图片
	imageURL, err := h.recognizeService.UploadImage(file)
	if err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	// 调用识别
	result, err := h.recognizeService.Recognize(userID, imageURL)
	if err != nil {
		response.Error(c, 500, err.Error(), nil)
		return
	}

	response.Success(c, "识别成功", result)
}

// GetHistory 获取识别历史
func (h *RecognizeHandler) GetHistory(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, 401, "未登录", nil)
		return
	}

	var query service.HistoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	records, total, err := h.recognizeService.GetHistory(userID, &query)
	if err != nil {
		response.Error(c, 500, err.Error(), nil)
		return
	}

	response.Success(c, "查询成功", gin.H{
		"list":      records,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	})
}

// DeleteHistory 删除识别记录
func (h *RecognizeHandler) DeleteHistory(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, 401, "未登录", nil)
		return
	}

	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	if err := h.recognizeService.DeleteHistory(userID, req.ID); err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "删除成功", nil)
}
