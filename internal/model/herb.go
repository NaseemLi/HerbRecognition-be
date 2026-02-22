package model

import (
	"time"

	"gorm.io/gorm"
)

// Herb 药材知识库表
type Herb struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:64;not null;uniqueIndex" json:"name"`
	Alias       string         `gorm:"size:255" json:"alias"`
	Description string         `gorm:"type:text" json:"description"`
	Effects     string         `gorm:"type:text" json:"effects"`
	Usage       string         `gorm:"type:text" json:"usage"`
	ImageURL    string         `gorm:"size:255" json:"image_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Herb) TableName() string {
	return "herbs"
}
