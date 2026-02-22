package admin

import (
	"herb-recognition-be/internal/service"
	"herb-recognition-be/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminUserHandler 管理端用户处理器
type AdminUserHandler struct {
	adminService *service.AdminHerbService
}

// NewAdminUserHandler 创建处理器
func NewAdminUserHandler() *AdminUserHandler {
	return &AdminUserHandler{
		adminService: &service.AdminHerbService{},
	}
}

// GetUserList 获取用户列表
func (h *AdminUserHandler) GetUserList(c *gin.Context) {
	var query service.UserListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	users, total, err := h.adminService.GetUserList(&query)
	if err != nil {
		response.Error(c, 500, err.Error(), nil)
		return
	}

	response.Success(c, "查询成功", gin.H{
		"list":      users,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	})
}

// UpdateUserRole 更新用户角色
func (h *AdminUserHandler) UpdateUserRole(c *gin.Context) {
	var req service.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	if err := h.adminService.UpdateUserRole(&req); err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "更新成功", nil)
}
