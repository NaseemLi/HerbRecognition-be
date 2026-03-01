package middleware

import (
	"net/http"
	"strings"

	"herb-recognition-be/pkg/jwtutil"
	"herb-recognition-be/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			response.Error(c, http.StatusUnauthorized, "未提供 token", nil)
			c.Abort()
			return
		}

		// 去除 Bearer 前缀
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		// 解析 token
		claims, err := jwtutil.ParseToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "token 无效或已过期", nil)
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RequireRole 角色权限校验中间件
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusForbidden, "未登录", nil)
			c.Abort()
			return
		}

		// admin 可以访问所有接口，user 只能访问用户接口
		if role == "user" && requiredRole == "admin" {
			response.Error(c, http.StatusForbidden, "权限不足", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserIDFromContext 从上下文获取用户 ID
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

// GetUsernameFromContext 从上下文获取用户名
func GetUsernameFromContext(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}
