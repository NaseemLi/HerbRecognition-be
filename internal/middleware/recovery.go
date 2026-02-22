package middleware

import (
	"herb-recognition-be/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery  Panic 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic recovered: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器内部错误",
					"data":    nil,
				})
			}
		}()
		c.Next()
	}
}
