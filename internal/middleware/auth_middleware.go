package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"movie-rating-api/internal/config"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	config *config.Config
}

// NewAuthMiddleware 创建认证中间件实例
func NewAuthMiddleware(config *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: config,
	}
}

// RequireAuth 要求认证的中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// 验证token
		token := parts[1]
		if token != m.config.AuthToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 认证成功，继续处理请求
		c.Next()
	}
}

// OptionalAuth 可选认证的中间件（用于健康检查等无需认证的端点）
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以添加日志记录等，但不阻止请求继续
		c.Next()
	}
}
