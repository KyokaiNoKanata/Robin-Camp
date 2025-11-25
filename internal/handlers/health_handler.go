package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {}

// NewHealthHandler 创建健康检查处理器实例
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check 健康检查
func (h *HealthHandler) Check(c *gin.Context) {
	// 健康检查端点，返回200状态码和OK消息
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "Movie Rating API is running",
	})
}
