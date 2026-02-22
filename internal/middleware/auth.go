package middleware

import (
	"net/http"
	"strings"

	"herb-recognition-be/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("herb-recognition-secret-key")

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
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "token 无效或已过期", nil)
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		claims, ok := token.Claims.(*Claims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "token 解析失败", nil)
			c.Abort()
			return
		}

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

		if role != requiredRole && requiredRole != "user" {
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
