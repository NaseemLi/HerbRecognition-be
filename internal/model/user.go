package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:32;not null" json:"username"`
	Password  string         `gorm:"size:64;not null" json:"-"`
	Role      string         `gorm:"size:20;default:user" json:"role"` // admin 或 user
	Avatar    string         `gorm:"size:255" json:"avatar"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
