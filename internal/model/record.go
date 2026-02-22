package model

import (
	"time"

	"gorm.io/gorm"
)

// RecognitionRecord 识别历史记录表
type RecognitionRecord struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	ImageURL   string         `gorm:"size:255;not null" json:"image_url"`
	HerbID     *uint          `gorm:"index" json:"herb_id"`
	HerbName   string         `gorm:"size:64" json:"herb_name"`
	Confidence float32        `gorm:"type:float(5,2)" json:"confidence"` // 置信度 0-100
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RecognitionRecord) TableName() string {
	return "recognition_records"
}
