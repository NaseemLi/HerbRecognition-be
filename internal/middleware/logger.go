package middleware

import (
	"herb-recognition-be/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)

		logger.Infof(
			"request method=%s path=%s status=%d latency=%s query=%q",
			c.Request.Method,
			path,
			c.Writer.Status(),
			cost,
			query,
		)
	}
}
