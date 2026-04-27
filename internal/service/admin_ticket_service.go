package service

import (
	"errors"
	"time"

	"herb-recognition-be/internal/model"
	"herb-recognition-be/internal/repository"

	"gorm.io/gorm"
)

// AdminTicketService 管理端工单服务
type AdminTicketService struct{}

// AdminTicketListQuery 管理端工单列表查询参数
type AdminTicketListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status"`
}

// TicketWithUser 带用户信息的工单
type TicketWithUser struct {
	model.Ticket
	Username string `json:"username"`
}

// GetTicketList 获取所有工单列表
func (s *AdminTicketService) GetTicketList(query *AdminTicketListQuery) ([]TicketWithUser, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 50 {
		query.PageSize = 10
	}

	offset := (query.Page - 1) * query.PageSize
	var tickets []TicketWithUser
	var total int64

	db := repository.DB.Model(&model.Ticket{}).
		Select("tickets.*, users.username").
		Joins("LEFT JOIN users ON tickets.user_id = users.id").
		Where("tickets.deleted_at IS NULL")

	if query.Status != "" {
		if !validTicketStatuses[query.Status] {
			return nil, 0, errors.New("无效的状态参数")
		}
		db = db.Where("tickets.status = ?", query.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, errors.New("查询失败")
	}

	if err := db.Order("tickets.created_at DESC").Offset(offset).Limit(query.PageSize).Find(&tickets).Error; err != nil {
		return nil, 0, errors.New("查询失败")
	}

	return tickets, total, nil
}

// GetTicketDetail 获取工单详情（管理端）
func (s *AdminTicketService) GetTicketDetail(ticketID uint) (*TicketWithUser, error) {
	var ticket TicketWithUser
	if err := repository.DB.Model(&model.Ticket{}).
		Select("tickets.*, users.username").
		Joins("LEFT JOIN users ON tickets.user_id = users.id").
		Where("tickets.id = ?", ticketID).
		First(&ticket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("工单不存在")
		}
		return nil, errors.New("查询失败")
	}
	return &ticket, nil
}

// ReplyTicketRequest 回复工单请求
type ReplyTicketRequest struct {
	Reply string `json:"reply" binding:"required,min=1"`
}

// ReplyTicket 回复工单
func (s *AdminTicketService) ReplyTicket(ticketID uint, req *ReplyTicketRequest) error {
	var ticket model.Ticket
	if err := repository.DB.First(&ticket, ticketID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("工单不存在")
		}
		return errors.New("查询失败")
	}

	if ticket.Status == "resolved" && ticket.AdminReply != "" {
		return errors.New("工单已回复，无法重复回复")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"admin_reply": req.Reply,
		"status":      "resolved",
		"replied_at":  &now,
	}

	if err := repository.DB.Model(&ticket).Updates(updates).Error; err != nil {
		return errors.New("回复失败")
	}

	return nil
}

// UpdateTicketStatusRequest 更新工单状态请求
type UpdateTicketStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending processing resolved"`
}

// UpdateTicketStatus 更新工单状态
func (s *AdminTicketService) UpdateTicketStatus(ticketID uint, req *UpdateTicketStatusRequest) error {
	var ticket model.Ticket
	if err := repository.DB.First(&ticket, ticketID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("工单不存在")
		}
		return errors.New("查询失败")
	}

	if ticket.Status == req.Status {
		return errors.New("工单状态已经是该状态")
	}

	if err := repository.DB.Model(&ticket).Update("status", req.Status).Error; err != nil {
		return errors.New("更新失败")
	}

	return nil
}
