package admin

import (
	"strconv"

	"herb-recognition-be/internal/service"
	"herb-recognition-be/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminTicketHandler 管理端工单处理器
type AdminTicketHandler struct {
	adminTicketService *service.AdminTicketService
}

// NewAdminTicketHandler 创建管理端工单处理器
func NewAdminTicketHandler() *AdminTicketHandler {
	return &AdminTicketHandler{
		adminTicketService: &service.AdminTicketService{},
	}
}

// GetTicketList 获取工单列表
func (h *AdminTicketHandler) GetTicketList(c *gin.Context) {
	var query service.AdminTicketListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	tickets, total, err := h.adminTicketService.GetTicketList(&query)
	if err != nil {
		response.Error(c, 500, err.Error(), nil)
		return
	}

	response.Success(c, "查询成功", gin.H{
		"list":      tickets,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	})
}

// GetTicketDetail 获取工单详情
func (h *AdminTicketHandler) GetTicketDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "参数错误", nil)
		return
	}

	ticket, err := h.adminTicketService.GetTicketDetail(uint(id))
	if err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "查询成功", ticket)
}

// ReplyTicket 回复工单
func (h *AdminTicketHandler) ReplyTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "参数错误", nil)
		return
	}

	var req service.ReplyTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	if err := h.adminTicketService.ReplyTicket(uint(id), &req); err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "回复成功", nil)
}

// UpdateTicketStatus 更新工单状态
func (h *AdminTicketHandler) UpdateTicketStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "参数错误", nil)
		return
	}

	var req service.UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	if err := h.adminTicketService.UpdateTicketStatus(uint(id), &req); err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "更新成功", nil)
}
