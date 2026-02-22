package model

import (
	"time"

	"gorm.io/gorm"
)

// Herb 中草药表
type Herb struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	Scientific string         `gorm:"size:200" json:"scientific"`
	Category   string         `gorm:"size:50" json:"category"`
	Effect     string         `gorm:"type:text" json:"effect"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// RecognitionRecord 识别记录表
type RecognitionRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ImagePath  string    `gorm:"size:500;not null" json:"image_path"`
	HerbID     uint      `gorm:"not null" json:"herb_id"`
	Confidence float32   `gorm:"type:decimal(5,4)" json:"confidence"`
	CreatedAt  time.Time `json:"created_at"`

	Herb Herb `gorm:"foreignKey:HerbID" json:"herb_detail"`
}
