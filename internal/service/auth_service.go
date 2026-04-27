package service

import (
	"errors"
	"mime/multipart"

	"herb-recognition-be/internal/model"
	"herb-recognition-be/internal/repository"
	"herb-recognition-be/pkg/jwtutil"
	"herb-recognition-be/pkg/upload"

	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务
type AuthService struct{}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) error {
	// 检查用户名是否已存在
	var existingUser model.User
	if err := repository.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return errors.New("用户名已存在")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     "user",
	}

	return repository.DB.Create(&user).Error
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	var user model.User
	if err := repository.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成 JWT Token
	token, err := jwtutil.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("token 生成失败")
	}

	return &LoginResponse{
		Token: token,
		User:  &user,
	}, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	// 查找用户
	var user model.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	return repository.DB.Model(&user).Update("password", string(hashedPassword)).Error
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=32"`
	Avatar   string `json:"avatar"`
}

// UpdateProfile 更新用户资料
func (s *AuthService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*model.User, error) {
	var user model.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 如果更新用户名，检查是否已存在
	if req.Username != "" && req.Username != user.Username {
		var existingUser model.User
		if err := repository.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
			return nil, errors.New("用户名已存在")
		}
		user.Username = req.Username
	}

	updates := map[string]interface{}{
		"avatar": req.Avatar,
	}

	if req.Username != "" && req.Username != user.Username {
		var existingUser model.User
		if err := repository.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
			return nil, errors.New("用户名已存在")
		}
		updates["username"] = req.Username
	}

	if len(updates) > 0 {
		if err := repository.DB.Model(&user).Updates(updates).Error; err != nil {
			return nil, errors.New("更新失败")
		}
	}

	// 同步更新的值到返回结构体
	if req.Username != "" {
		user.Username = req.Username
	}
	user.Avatar = req.Avatar

	return &user, nil
}

// UploadAvatar 上传头像
func (s *AuthService) UploadAvatar(file *multipart.FileHeader) (string, error) {
	cfg := upload.AvatarConfig
	cfg.UploadDir = "./uploads/avatars"
	cfg.URLPrefix = "/uploads/avatars/"
	return upload.UploadFile(file, cfg)
}

// GetProfile 获取用户资料
func (s *AuthService) GetProfile(userID uint) (*model.User, error) {
	var user model.User
	if err := repository.DB.Select("id, username, role, avatar, created_at").First(&user, userID).Error; err != nil {
		return nil, errors.New("获取用户资料失败")
	}
	return &user, nil
}
