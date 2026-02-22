package app

import (
	"herb-recognition-be/internal/middleware"
	"herb-recognition-be/internal/service"
	"herb-recognition-be/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: &service.AuthService{},
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	if err := h.authService.Register(&req); err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "注册成功", nil)
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	loginResp, err := h.authService.Login(&req)
	if err != nil {
		response.Error(c, 401, err.Error(), nil)
		return
	}

	response.Success(c, "登录成功", loginResp)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, 401, "未登录", nil)
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	if err := h.authService.ChangePassword(userID, &req); err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "密码修改成功", nil)
}
