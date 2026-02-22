package admin

import (
	"herb-recognition-be/internal/service"
	"herb-recognition-be/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminHerbHandler 管理端药材处理器
type AdminHerbHandler struct {
	adminHerbService *service.AdminHerbService
}

// NewAdminHerbHandler 创建处理器
func NewAdminHerbHandler() *AdminHerbHandler {
	return &AdminHerbHandler{
		adminHerbService: &service.AdminHerbService{},
	}
}

// CreateHerb 新增药材
func (h *AdminHerbHandler) CreateHerb(c *gin.Context) {
	var req service.CreateHerbRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	herb, err := h.adminHerbService.CreateHerb(&req)
	if err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "创建成功", herb)
}

// UpdateHerb 编辑药材
func (h *AdminHerbHandler) UpdateHerb(c *gin.Context) {
	var req service.UpdateHerbRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	herb, err := h.adminHerbService.UpdateHerb(&req)
	if err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "更新成功", herb)
}

// DeleteHerb 删除药材
func (h *AdminHerbHandler) DeleteHerb(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	if err := h.adminHerbService.DeleteHerb(req.ID); err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "删除成功", nil)
}

// BatchDelete 批量删除药材
func (h *AdminHerbHandler) BatchDelete(c *gin.Context) {
	var req struct {
		Ids []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	if err := h.adminHerbService.BatchDeleteHerb(req.Ids); err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "删除成功", nil)
}

// UploadAndSetImage 上传并设置药材图片
func (h *AdminHerbHandler) UploadAndSetImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		response.Error(c, 400, "请上传图片文件", nil)
		return
	}

	imageURL, err := h.adminHerbService.UploadAndSetImage(file)
	if err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "上传成功", gin.H{"image_url": imageURL})
}

// GetList 获取药材列表
func (h *AdminHerbHandler) GetList(c *gin.Context) {
	var query service.HerbListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	herbs, total, err := h.adminHerbService.GetHerbList(&query)
	if err != nil {
		response.Error(c, 500, err.Error(), nil)
		return
	}

	response.Success(c, "查询成功", gin.H{
		"list":      herbs,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	})
}
