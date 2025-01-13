package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"

	"hy2agent/internal/config"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()
		
		// 获取客户端Host
		clientHost := c.Request.Host
		if host, _, err := net.SplitHostPort(clientHost); err == nil {
			clientHost = host
		}
		
		// 检查黑名单
		for _, ip := range cfg.IPBlacklist {
			if ip == clientIP {
				c.JSON(403, gin.H{"error": "IP is blacklisted"})
				c.Abort()
				return
			}
		}
		
		// 检查白名单
		allowed := false
		
		// 检查IP白名单
		for _, ip := range cfg.IPWhitelist {
			if ip == clientIP {
				allowed = true
				break
			}
		}
		
		// 检查域名白名单
		if !allowed && len(cfg.DomainWhitelist) > 0 {
			for _, domain := range cfg.DomainWhitelist {
				if domain == clientHost {
					allowed = true
					break
				}
				// 支持通配符域名
				if strings.HasPrefix(domain, "*.") && 
				   strings.HasSuffix(clientHost, domain[1:]) {
					allowed = true
					break
				}
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