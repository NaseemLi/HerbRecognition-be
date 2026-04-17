package main

import (
	"crypto/rand"
	"encoding/hex"
	"herb-recognition-be/internal/config"
	"herb-recognition-be/internal/middleware"
	"herb-recognition-be/internal/model"
	"herb-recognition-be/internal/repository"
	"herb-recognition-be/internal/routes"
	"herb-recognition-be/pkg/logger"
	"herb-recognition-be/pkg/onnx"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// generateRandomPassword 生成随机密码
func generateRandomPassword(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "changeme123" // fallback
	}
	return hex.EncodeToString(bytes)[:length]
}

// initRootUser 创建默认 root 管理员用户
func initRootUser() {
	username := config.Conf.Admin.Username
	if username == "" {
		username = "root"
	}

	var user model.User
	// 检查是否已存在 root 用户
	if err := repository.DB.Where("username = ?", username).First(&user).Error; err == nil {
		logger.Infof("%s 用户已存在，跳过创建", username)
		return
	}

	// 从配置获取密码，如果为空则生成随机密码
	password := config.Conf.Admin.Password
	isRandomPassword := false
	if password == "" {
		password = generateRandomPassword(12)
		isRandomPassword = true
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("%s 用户密码加密失败：%v", username, err)
		return
	}

	// 创建 root 用户
	rootUser := model.User{
		Username: username,
		Password: string(hashedPassword),
		Role:     "admin",
		Avatar:   "",
	}

	if err := repository.DB.Create(&rootUser).Error; err != nil {
		logger.Errorf("创建 %s 用户失败：%v", username, err)
		return
	}

	if isRandomPassword {
		logger.Infof("%s 用户创建成功，随机密码：%s（请登录后立即修改密码）", username, password)
	} else {
		logger.Infof("%s 用户创建成功（请登录后修改密码）", username)
	}
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

	// 初始化 ONNX 预测器
	modelPath := config.Conf.ModelService.ONNXModelPath
	if modelPath == "" {
		modelPath = "./models/onnx/herb.onnx"
	}
	classesPath := config.Conf.ModelService.ClassesPath
	if classesPath == "" {
		classesPath = "./models/onnx/classes.txt"
	}

	if err := onnx.InitPredictor(modelPath, classesPath); err != nil {
		logger.Warnf("ONNX 预测器初始化失败：%v", err)
		logger.Warn("识别功能将不可用，请检查模型文件")
	} else {
		logger.Infof("ONNX 预测器初始化成功，类别数：%d", onnx.GetClassCount())
	}

	serverMode := config.Conf.Server.Mode
	if serverMode == "" {
		serverMode = gin.ReleaseMode
	}
	gin.SetMode(serverMode)

	// 创建 Gin 路由
	r := gin.New()

	// 设置可信代理
	r.SetTrustedProxies(nil)

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
