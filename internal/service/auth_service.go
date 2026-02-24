package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"herb-recognition-be/internal/model"
	"herb-recognition-be/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

// jwtSecret JWT 密钥
var jwtSecret = []byte("herb-recognition-secret-key")

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
	token, err := generateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("token 生成失败")
	}

	return &LoginResponse{
		Token: token,
		User:  &user,
	}, nil
}

// generateToken 生成 JWT Token
func generateToken(userID uint, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 天有效期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
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

	// 更新头像
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	// 保存更新
	if err := repository.DB.Save(&user).Error; err != nil {
		return nil, errors.New("更新失败")
	}

	return &user, nil
}

// UploadAvatar 上传头像
func (s *AuthService) UploadAvatar(file *multipart.FileHeader) (string, error) {
	// 检查文件大小
	if file.Size > 2*1024*1024 { // 2MB 限制
		return "", errors.New("头像文件大小不能超过 2MB")
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExtensions[ext] {
		return "", errors.New("不支持的头像格式，仅支持 JPG、PNG、GIF、WEBP")
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", errors.New("文件打开失败")
	}
	defer src.Close()

	// 读取文件头进行 MIME 类型校验
	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && err != io.EOF {
		return "", errors.New("文件读取失败")
	}
	buf = buf[:n]

	// 检测真实 MIME 类型
	contentType := http.DetectContentType(buf)
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedMimeTypes[contentType] {
		return "", errors.New("文件类型不匹配，请上传有效的图片")
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("avatar_%s%s", uuid.New().String(), ext)
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", errors.New("创建上传目录失败")
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", errors.New("文件保存失败")
	}
	defer dst.Close()

	// 重置文件指针并复制内容
	src.Seek(0, 0)
	if _, err := io.Copy(dst, src); err != nil {
		return "", errors.New("文件保存失败")
	}

	// 返回访问 URL
	url := "/uploads/avatars/" + filename
	return url, nil
}
