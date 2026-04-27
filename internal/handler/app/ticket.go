package app

import (
	"strconv"

	"herb-recognition-be/internal/middleware"
	"herb-recognition-be/internal/service"
	"herb-recognition-be/pkg/response"

	"github.com/gin-gonic/gin"
)

// TicketHandler 工单处理器
type TicketHandler struct {
	ticketService *service.TicketService
}

// NewTicketHandler 创建工单处理器
func NewTicketHandler() *TicketHandler {
	return &TicketHandler{
		ticketService: &service.TicketService{},
	}
}

// CreateTicket 提交工单
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		response.Error(c, 401, "未登录", nil)
		return
	}

	var req service.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	ticket, err := h.ticketService.CreateTicket(userID, &req)
	if err != nil {
		response.Error(c, 400, err.Error(), nil)
		return
	}

	response.Success(c, "提交成功", ticket)
}

// GetMyTickets 获取我的工单列表
func (h *TicketHandler) GetMyTickets(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		response.Error(c, 401, "未登录", nil)
		return
	}

	var query service.TicketListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, "请求参数错误", nil)
		return
	}

	tickets, total, err := h.ticketService.GetMyTickets(userID, &query)
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
func (h *TicketHandler) GetTicketDetail(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		response.Error(c, 401, "未登录", nil)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "参数错误", nil)
		return
	}

	ticket, err := h.ticketService.GetTicketDetail(userID, uint(id))
	if err != nil {
		if err.Error() == "工单不存在" {
			response.Error(c, 404, err.Error(), nil)
			return
		}
		response.Error(c, 500, err.Error(), nil)
		return
	}

	response.Success(c, "查询成功", ticket)
}
