package repository

import (
	"fmt"
	"herb-recognition-be/internal/config"
	"herb-recognition-be/internal/model"
	"herb-recognition-be/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	// 从配置中获取 DSN 连接字符串
	dsn := config.Conf.Database.BuildDSN()

	// 打开数据库连接
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 自动迁移表结构（根据 model 创建/更新表）
	err = DB.AutoMigrate(&model.User{}, &model.Herb{}, &model.RecognitionRecord{})
	if err != nil {
		return fmt.Errorf("数据库自动迁移失败：%w", err)
	}

	logger.Info("数据库连接成功并已完成迁移")
	return nil
}
