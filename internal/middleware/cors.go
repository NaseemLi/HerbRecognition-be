package middleware

import (
	"herb-recognition-be/internal/config"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从配置获取允许的域名
		allowOrigins := config.Conf.CORS.AllowOrigins
		origin := c.GetHeader("Origin")

		// 检查是否在白名单中
		allowedOrigin := ""
		for _, o := range allowOrigins {
			if o == origin {
				allowedOrigin = origin
				break
			}
		}

		// 如果在白名单中，设置对应的 origin
		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		} else if len(allowOrigins) == 0 {
			// 如果没有配置白名单，开发模式下允许所有来源
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
