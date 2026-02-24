package main

import (
	"herb-recognition-be/internal/config"
	"herb-recognition-be/internal/middleware"
	"herb-recognition-be/internal/model"
	"herb-recognition-be/internal/repository"
	"herb-recognition-be/internal/routes"
	"herb-recognition-be/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// initRootUser 创建默认 root 管理员用户
func initRootUser() {
	var user model.User
	// 检查是否已存在 root 用户
	if err := repository.DB.Where("username = ?", "root").First(&user).Error; err == nil {
		logger.Info("root 用户已存在，跳过创建")
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("root 用户密码加密失败：%v", err)
		return
	}

	// 创建 root 用户
	rootUser := model.User{
		Username: "root",
		Password: string(hashedPassword),
		Role:     "admin",
		Avatar:   "",
	}

	if err := repository.DB.Create(&rootUser).Error; err != nil {
		logger.Error("创建 root 用户失败：%v", err)
		return
	}

	logger.Info("root 用户创建成功（用户名：root, 密码：admin123）")
}

func main() {
	// 初始化日志（默认 info 级别）
	if err := logger.Init("info"); err != nil {
		logger.Fatalf("日志初始化失败：%v", err)
	}

	// 加载配置
	if err := config.Init(); err != nil {
		logger.Fatalf("配置加载失败：%v", err)
	}
	logger.Info("配置加载成功")

	// 重新初始化日志（使用配置中的级别）
	if err := logger.Init(config.Conf.Server.Mode); err != nil {
		logger.Fatalf("日志初始化失败：%v", err)
	}
	logger.Info("日志初始化成功")

	// 初始化数据库
	if err := repository.InitDB(); err != nil {
		logger.Fatalf("数据库初始化失败：%v", err)
	}
	logger.Info("数据库初始化成功")

	// 创建默认 root 用户
	initRootUser()

	// 创建 Gin 路由
	r := gin.New()

	// 设置可信代理
	r.SetTrustedProxies(nil)

	// 设置 Gin 模式
	gin.SetMode(config.Conf.Server.Mode)

	// 注册全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 注册路由
	routes.InitRoutes(r)

	// 启动服务
	port := config.Conf.Server.Port

	if err := r.Run(":" + port); err != nil {
		logger.Fatalf("服务启动失败：%v", err)
	}
}
