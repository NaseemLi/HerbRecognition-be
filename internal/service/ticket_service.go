package service

import (
	"errors"

	"herb-recognition-be/internal/model"
	"herb-recognition-be/internal/repository"

	"gorm.io/gorm"
)

// TicketService 工单服务
type TicketService struct{}

var validTicketStatuses = map[string]bool{
	"pending":    true,
	"processing": true,
	"resolved":   true,
}

// CreateTicketRequest 创建工单请求
type CreateTicketRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=128"`
	Content  string `json:"content" binding:"required,min=1"`
	ImageURL string `json:"image_url" binding:"omitempty,max=255"`
}

// CreateTicket 创建工单
func (s *TicketService) CreateTicket(userID uint, req *CreateTicketRequest) (*model.Ticket, error) {
	ticket := model.Ticket{
		UserID:   userID,
		Title:    req.Title,
		Content:  req.Content,
		ImageURL: req.ImageURL,
		Status:   "pending",
	}

	if err := repository.DB.Create(&ticket).Error; err != nil {
		return nil, errors.New("提交工单失败")
	}

	return &ticket, nil
}

// TicketListQuery 工单列表查询参数
type TicketListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status"`
}

// GetMyTickets 获取我的工单列表
func (s *TicketService) GetMyTickets(userID uint, query *TicketListQuery) ([]model.Ticket, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 50 {
		query.PageSize = 10
	}

	offset := (query.Page - 1) * query.PageSize
	var tickets []model.Ticket
	var total int64

	db := repository.DB.Model(&model.Ticket{}).Where("user_id = ?", userID)
	if query.Status != "" {
		if !validTicketStatuses[query.Status] {
			return nil, 0, errors.New("无效的状态参数")
		}
		db = db.Where("status = ?", query.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, errors.New("查询失败")
	}

	if err := db.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&tickets).Error; err != nil {
		return nil, 0, errors.New("查询失败")
	}

	return tickets, total, nil
}

// GetTicketDetail 获取工单详情
func (s *TicketService) GetTicketDetail(userID, ticketID uint) (*model.Ticket, error) {
	var ticket model.Ticket
	if err := repository.DB.Where("id = ? AND user_id = ?", ticketID, userID).First(&ticket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("工单不存在")
		}
		return nil, errors.New("查询失败")
	}
	return &ticket, nil
}
