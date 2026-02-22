package main

import (
	"fmt"
	"log"

	"herb-recognition-be/internal/config"
	"herb-recognition-be/internal/middleware"
	"herb-recognition-be/internal/repository"
	"herb-recognition-be/internal/routes"
	"herb-recognition-be/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.Init(); err != nil {
		log.Fatalf("配置加载失败：%v", err)
	}
	fmt.Println("配置加载成功!")

	// 初始化日志
	if err := logger.Init(config.Conf.Server.Mode); err != nil {
		log.Fatalf("日志初始化失败：%v", err)
	}
	fmt.Println("日志初始化成功!")

	// 初始化数据库
	if err := repository.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败：%v", err)
	}
	fmt.Println("数据库初始化成功!")

	// 创建 Gin 路由
	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())
	r.Use(gin.Logger())

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 注册路由
	routes.InitRoutes(r)

	//  启动服务
	port := config.Conf.Server.Port
	fmt.Printf("服务启动在 http://localhost:%s\n", port)
	logger.Infof("服务启动在端口 %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务启动失败：%v", err)
	}
}
