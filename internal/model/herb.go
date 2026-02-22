package model

import (
	"time"

	"gorm.io/gorm"
)

// Herb 药材知识库表
type Herb struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:64;not null;index:idx_name" json:"name"` // 药材名
	Scientific  string         `gorm:"size:100" json:"scientific"`                  // 学名
	Alias       string         `gorm:"size:255" json:"alias"`                       // 别名（逗号分隔）
	Category    string         `gorm:"size:32;index:idx_category" json:"category"`  // 分类
	Description string         `gorm:"type:text" json:"description"`                // 描述
	Effects     string         `gorm:"type:text" json:"effects"`                    // 功效
	Usage       string         `gorm:"type:text" json:"usage"`                      // 用法用量
	ImageURL    string         `gorm:"size:255" json:"image_url"`                   // 图片
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Herb) TableName() string {
	return "herbs"
}
