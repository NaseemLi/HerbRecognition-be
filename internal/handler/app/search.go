package app

import (
	"herb-recognition-be/internal/service"
	"herb-recognition-be/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SearchHandler 搜索处理器
type SearchHandler struct {
	herbService *service.HerbService
}

// NewSearchHandler 创建搜索处理器
func NewSearchHandler() *SearchHandler {
	return &SearchHandler{
		herbService: &service.HerbService{},
	}
}

// Search 搜索药材
func (h *SearchHandler) Search(c *gin.Context) {
	var req service.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	result, err := h.herbService.Search(&req)
	if err != nil {
		response.Error(c, 500, err.Error(), nil)
		return
	}

	response.Success(c, "查询成功", result)
}

// GetDetail 获取药材详情
func (h *SearchHandler) GetDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的 ID", nil)
		return
	}

	herb, err := h.herbService.GetDetail(uint(id))
	if err != nil {
		response.Error(c, 404, err.Error(), nil)
		return
	}

	response.Success(c, "查询成功", herb)
}

// GetAll 获取药材列表
func (h *SearchHandler) GetAll(c *gin.Context) {
	var req service.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	result, err := h.herbService.GetAll(&req)
	if err != nil {
		response.Error(c, 500, err.Error(), nil)
		return
	}

	response.Success(c, "查询成功", result)
}
