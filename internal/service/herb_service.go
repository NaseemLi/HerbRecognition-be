package service

import (
	"errors"
	"herb-recognition-be/internal/model"
	"herb-recognition-be/internal/repository"
	"strings"

	"gorm.io/gorm"
)

// HerbService 药材服务
type HerbService struct{}

// SearchRequest 搜索请求
type SearchRequest struct {
	Keyword  string `form:"keyword" binding:"required,min=1,max=32"` // 搜索关键词
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	List     []model.Herb `json:"list"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

// Search 搜索药材（模糊搜索）
func (s *HerbService) Search(req *SearchRequest) (*SearchResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 50 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize
	keyword := "%" + strings.TrimSpace(req.Keyword) + "%"

	var herbs []model.Herb
	var total int64

	// 查询总数
	if err := repository.DB.Model(&model.Herb{}).
		Where("name LIKE ? OR alias LIKE ? OR scientific LIKE ?", keyword, keyword, keyword).
		Count(&total).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	// 分页查询
	if err := repository.DB.
		Where("name LIKE ? OR alias LIKE ? OR scientific LIKE ?", keyword, keyword, keyword).
		Offset(offset).Limit(req.PageSize).
		Find(&herbs).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	return &SearchResponse{
		List:     herbs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetDetail 获取药材详情
func (s *HerbService) GetDetail(id uint) (*model.Herb, error) {
	var herb model.Herb
	if err := repository.DB.First(&herb, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("药材不存在")
		}
		return nil, errors.New("查询失败")
	}

	return &herb, nil
}

// GetByName 根据名称获取药材
func (s *HerbService) GetByName(name string) (*model.Herb, error) {
	var herb model.Herb
	if err := repository.DB.Where("name = ?", name).First(&herb).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("药材不存在")
		}
		return nil, errors.New("查询失败")
	}

	return &herb, nil
}

// GetAll 获取所有药材（分页）
type GetAllRequest struct {
	Category string `form:"category"` // 按分类筛选
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

func (s *HerbService) GetAll(req *GetAllRequest) (*SearchResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 50 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize

	var herbs []model.Herb
	var total int64

	query := repository.DB.Model(&model.Herb{})

	// 按分类筛选
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	// 分页查询
	if err := query.Offset(offset).Limit(req.PageSize).Find(&herbs).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	return &SearchResponse{
		List:     herbs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
