package app

import (
	"herb-recognition-be/internal/middleware"
	"herb-recognition-be/internal/service"
	"herb-recognition-be/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	authService *service.AuthService
}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		authService: &service.AuthService{},
	}
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, 401, "未登录", nil)
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	user, err := h.authService.UpdateProfile(userID, &req)
	if err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "资料更新成功", gin.H{
		"user": user,
	})
}

// UploadAvatar 上传头像
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, 401, "未登录", nil)
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		response.Error(c, 400, "请选择头像文件", nil)
		return
	}

	avatarURL, err := h.authService.UploadAvatar(file)
	if err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	// 同时更新用户资料中的头像字段
	updateReq := &service.UpdateProfileRequest{
		Avatar: avatarURL,
	}
	_, err = h.authService.UpdateProfile(userID, updateReq)
	if err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "头像上传成功", gin.H{
		"avatar_url": avatarURL,
	})
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, 401, "未登录", nil)
		return
	}

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "获取成功", gin.H{
		"user": user,
	})
}
