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

	"herb-recognition-be/internal/model"
	"herb-recognition-be/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AdminHerbService 管理端药材服务
type AdminHerbService struct{}

// CreateHerbRequest 创建药材请求
type CreateHerbRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=64"`
	Scientific  string `json:"scientific"`
	Alias       string `json:"alias"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Effects     string `json:"effects"`
	Usage       string `json:"usage"`
}

// CreateHerb 创建药材
func (s *AdminHerbService) CreateHerb(req *CreateHerbRequest) (*model.Herb, error) {
	// 检查名称是否已存在
	var existing model.Herb
	if err := repository.DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		return nil, errors.New("药材名称已存在")
	}

	herb := model.Herb{
		Name:        req.Name,
		Scientific:  req.Scientific,
		Alias:       req.Alias,
		Category:    req.Category,
		Description: req.Description,
		Effects:     req.Effects,
		Usage:       req.Usage,
	}

	if err := repository.DB.Create(&herb).Error; err != nil {
		return nil, errors.New("创建失败")
	}

	return &herb, nil
}

// UpdateHerbRequest 更新药材请求
type UpdateHerbRequest struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required,min=1,max=64"`
	Scientific  string `json:"scientific"`
	Alias       string `json:"alias"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Effects     string `json:"effects"`
	Usage       string `json:"usage"`
}

// UpdateHerb 更新药材
func (s *AdminHerbService) UpdateHerb(req *UpdateHerbRequest) (*model.Herb, error) {
	var herb model.Herb
	if err := repository.DB.First(&herb, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("药材不存在")
		}
		return nil, errors.New("查询失败")
	}

	// 检查新名称是否与其他记录冲突
	var existing model.Herb
	if err := repository.DB.Where("name = ? AND id != ?", req.Name, req.ID).First(&existing).Error; err == nil {
		return nil, errors.New("药材名称已存在")
	}

	herb.Name = req.Name
	herb.Scientific = req.Scientific
	herb.Alias = req.Alias
	herb.Category = req.Category
	herb.Description = req.Description
	herb.Effects = req.Effects
	herb.Usage = req.Usage

	if err := repository.DB.Save(&herb).Error; err != nil {
		return nil, errors.New("更新失败")
	}

	return &herb, nil
}

// DeleteHerb 删除药材
func (s *AdminHerbService) DeleteHerb(id uint) error {
	var herb model.Herb
	if err := repository.DB.First(&herb, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("药材不存在")
		}
		return errors.New("查询失败")
	}

	// 删除关联的图片文件
	if herb.ImageURL != "" {
		filePath := "." + herb.ImageURL
		os.Remove(filePath)
	}

	// 删除药材（软删除）
	if err := repository.DB.Delete(&herb).Error; err != nil {
		return errors.New("删除失败")
	}

	return nil
}

// UploadAndSetImage 上传并设置药材图片
func (s *AdminHerbService) UploadAndSetImage(file *multipart.FileHeader) (string, error) {
	// 检查文件大小 (5MB)
	if file.Size > 5*1024*1024 {
		return "", errors.New("文件大小不能超过 5MB")
	}

	// 检查扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		return "", errors.New("不支持的图片格式")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", errors.New("文件打开失败")
	}
	defer src.Close()

	// 读取文件头校验 MIME 类型
	buf := make([]byte, 512)
	n, _ := src.Read(buf)
	contentType := http.DetectContentType(buf[:n])
	allowedMimes := map[string]bool{
		"image/jpeg": true, "image/png": true, "image/gif": true, "image/webp": true,
	}
	if !allowedMimes[contentType] {
		return "", errors.New("文件类型不匹配")
	}

	// 生成文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadDir := "./uploads/herbs"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", errors.New("创建目录失败")
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", errors.New("文件保存失败")
	}
	defer dst.Close()

	src.Seek(0, 0)
	if _, err := io.Copy(dst, src); err != nil {
		return "", errors.New("文件保存失败")
	}

	return "/uploads/herbs/" + filename, nil
}

// BatchDeleteHerb 批量删除药材
func (s *AdminHerbService) BatchDeleteHerb(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的药材")
	}

	// 删除关联的图片文件
	var herbs []model.Herb
	if err := repository.DB.Where("id IN ?", ids).Find(&herbs).Error; err != nil {
		return errors.New("查询失败")
	}

	for _, herb := range herbs {
		if herb.ImageURL != "" {
			filePath := "." + herb.ImageURL
			os.Remove(filePath)
		}
	}

	// 批量删除（软删除）
	if err := repository.DB.Where("id IN ?", ids).Delete(&model.Herb{}).Error; err != nil {
		return errors.New("删除失败")
	}

	return nil
}

// HerbListQuery 药材列表查询参数
type HerbListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Category string `form:"category"`
}

// GetHerbList 获取药材列表（分页）
func (s *AdminHerbService) GetHerbList(query *HerbListQuery) ([]model.Herb, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 50 {
		query.PageSize = 10
	}

	offset := (query.Page - 1) * query.PageSize
	var herbs []model.Herb
	var total int64

	db := repository.DB.Model(&model.Herb{})
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, errors.New("查询失败")
	}

	if err := db.Order("id DESC").Offset(offset).Limit(query.PageSize).Find(&herbs).Error; err != nil {
		return nil, 0, errors.New("查询失败")
	}

	return herbs, total, nil
}

// UserListQuery 用户列表查询参数
type UserListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Role     string `form:"role"` // 按角色筛选
}

// GetUserList 获取用户列表（分页）
func (s *AdminHerbService) GetUserList(query *UserListQuery) ([]model.User, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 50 {
		query.PageSize = 10
	}

	offset := (query.Page - 1) * query.PageSize
	var users []model.User
	var total int64

	db := repository.DB.Model(&model.User{})
	if query.Role != "" {
		db = db.Where("role = ?", query.Role)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, errors.New("查询失败")
	}

	if err := db.Order("id DESC").Offset(offset).Limit(query.PageSize).Find(&users).Error; err != nil {
		return nil, 0, errors.New("查询失败")
	}

	return users, total, nil
}

// UpdateUserRoleRequest 更新用户角色请求
type UpdateUserRoleRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=user admin"`
}

// UpdateUserRole 更新用户角色
func (s *AdminHerbService) UpdateUserRole(req *UpdateUserRoleRequest) error {
	var user model.User
	if err := repository.DB.First(&user, req.UserID).Error; err != nil {
		return errors.New("用户不存在")
	}

	user.Role = req.Role
	return repository.DB.Save(&user).Error
}
