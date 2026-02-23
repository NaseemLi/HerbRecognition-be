package model

import (
	"time"

	"gorm.io/gorm"
)

// Herb 药材知识库表
type Herb struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:64;not null;index:idx_name" json:"name"`
	Scientific  string         `gorm:"size:100" json:"scientific"`
	Alias       string         `gorm:"size:255" json:"alias"`
	Category    string         `gorm:"size:32;index:idx_category" json:"category"`
	Description string         `gorm:"type:text" json:"description"`
	Effects     string         `gorm:"type:text" json:"effects"`
	Usage       string         `gorm:"type:text" json:"usage"`
	ImageURL    string         `gorm:"size:255" json:"image_url"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Herb) TableName() string {
	return "herbs"
}
