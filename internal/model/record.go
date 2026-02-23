package model

import (
	"time"

	"gorm.io/gorm"
)

// RecognitionRecord 识别历史记录表
type RecognitionRecord struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"-"` // 不暴露给前端
	ImageURL   string         `gorm:"size:255;not null" json:"image_url"`
	HerbID     *uint          `gorm:"index" json:"herb_id"`
	HerbName   string         `gorm:"size:64" json:"herb_name"`
	Confidence float32        `gorm:"type:float(5,2)" json:"confidence"`
	Status     int            `gorm:"default:1" json:"-"` // 内部状态
	ErrMsg     string         `gorm:"size:255" json:"-"`  // 内部错误信息
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RecognitionRecord) TableName() string {
	return "recognition_records"
}
