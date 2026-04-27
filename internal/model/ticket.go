package model

import (
	"time"

	"gorm.io/gorm"
)

// Ticket 工单表
type Ticket struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	Title      string         `gorm:"size:128;not null" json:"title"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	ImageURL   string         `gorm:"size:255" json:"image_url"`
	Status     string         `gorm:"size:20;default:pending" json:"status"`
	AdminReply string         `gorm:"type:text" json:"admin_reply"`
	RepliedAt  *time.Time     `json:"replied_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Ticket) TableName() string {
	return "tickets"
}
