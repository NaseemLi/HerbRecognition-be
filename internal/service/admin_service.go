package service

import (
	"errors"
	"mime/multipart"

	"herb-recognition-be/internal/model"
	"herb-recognition-be/internal/repository"
	"herb-recognition-be/pkg/upload"

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
	ImageUrl    string `json:"image_url"`
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
		ImageURL:    req.ImageUrl,
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
	ImageUrl    string `json:"image_url"`
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

	updates := map[string]interface{}{
		"name":        req.Name,
		"scientific":  req.Scientific,
		"alias":       req.Alias,
		"category":    req.Category,
		"description": req.Description,
		"effects":     req.Effects,
		"usage":       req.Usage,
		"image_url":   req.ImageUrl,
	}

	if err := repository.DB.Model(&herb).Updates(updates).Error; err != nil {
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
	upload.DeleteFile(herb.ImageURL)

	// 删除药材（软删除）
	if err := repository.DB.Delete(&herb).Error; err != nil {
		return errors.New("删除失败")
	}

	return nil
}

// UploadAndSetImage 上传并设置药材图片
func (s *AdminHerbService) UploadAndSetImage(file *multipart.FileHeader) (string, error) {
	cfg := upload.DefaultImageConfig
	cfg.UploadDir = "./uploads/herbs"
	cfg.URLPrefix = "/uploads/herbs/"
	return upload.UploadFile(file, cfg)
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
		upload.DeleteFile(herb.ImageURL)
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

	return repository.DB.Model(&user).Update("role", req.Role).Error
}

// DeleteUser 删除用户（管理员不能删除自己和其他管理员）
func (s *AdminHerbService) DeleteUser(userID uint, adminID uint) error {
	// 不能删除自己
	if userID == adminID {
		return errors.New("不能删除自己")
	}

	var user model.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("用户不存在")
		}
		return errors.New("查询失败")
	}

	// 不能删除管理员
	if user.Role == "admin" {
		return errors.New("不能删除管理员")
	}

	if err := repository.DB.Delete(&user).Error; err != nil {
		return errors.New("删除失败")
	}

	return nil
}

// BatchDeleteUser 批量删除用户
func (s *AdminHerbService) BatchDeleteUser(userIDs []uint, adminID uint) error {
	if len(userIDs) == 0 {
		return errors.New("请选择要删除的用户")
	}

	// 去重
	seen := make(map[uint]struct{})
	unique := make([]uint, 0, len(userIDs))
	for _, id := range userIDs {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			unique = append(unique, id)
		}
	}

	// 检查是否包含自己
	for _, id := range unique {
		if id == adminID {
			return errors.New("不能删除自己")
		}
	}

	// 检查是否存在管理员角色
	var adminCount int64
	if err := repository.DB.Model(&model.User{}).Where("id IN ? AND role = ?", unique, "admin").Count(&adminCount).Error; err != nil {
		return errors.New("查询失败")
	}
	if adminCount > 0 {
		return errors.New("不能删除管理员")
	}

	// 检查目标用户是否全部存在
	var existingCount int64
	if err := repository.DB.Model(&model.User{}).Where("id IN ?", unique).Count(&existingCount).Error; err != nil {
		return errors.New("查询失败")
	}
	if existingCount != int64(len(unique)) {
		return errors.New("部分用户不存在")
	}

	if err := repository.DB.Where("id IN ?", unique).Delete(&model.User{}).Error; err != nil {
		return errors.New("删除失败")
	}

	return nil
}
