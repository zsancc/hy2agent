package middleware

import (
	"github.com/gin-gonic/gin"

	"hy2agent/internal/config"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()

		// 检查白名单
		allowed := false
		for _, ip := range cfg.IPWhitelist {
			if ip == clientIP {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(403, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		// API Key认证
		reqApiKey := c.GetHeader("X-API-Key")
		if reqApiKey == "" {
			c.JSON(401, gin.H{"error": "No API key provided"})
			c.Abort()
			return
		}

		if reqApiKey != cfg.APIKey {
			c.JSON(401, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
