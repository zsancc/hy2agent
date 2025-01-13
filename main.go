package main

import (
    "github.com/gin-gonic/gin"
    v1 "agentapi/api/v1"
    "agentapi/internal/config"
    "log"
)

func main() {
    // 加载配置
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // 打印API Key，仅在首次安装时显示
    log.Printf("Agent API Key: %s", cfg.APIKey)
    
    r := gin.Default()
    
    // API认证中间件
    r.Use(authMiddleware(cfg))
    
    // API路由
    statusHandler := v1.NewStatusHandler()
    systemHandler := v1.NewSystemHandler()
    hysteria2Handler := v1.NewHysteria2Handler()
    
    // 状态API
    r.GET("/api/v1/status", statusHandler.GetStatus)
    
    // 系统管理API
    systemGroup := r.Group("/api/v1/system")
    {
        systemGroup.GET("/memory", systemHandler.GetMemory)
        systemGroup.GET("/disk", systemHandler.GetDisk)
        systemGroup.GET("/network", systemHandler.GetNetwork)
        systemGroup.GET("/info", systemHandler.GetInfo)
        systemGroup.POST("/reboot", systemHandler.Reboot)
        systemGroup.POST("/shutdown", systemHandler.Shutdown)
    }
    
    // Hysteria2管理API
    hysteria2Group := r.Group("/api/v1/hysteria")
    {
        hysteria2Group.GET("/status", hysteria2Handler.GetStatus)
        hysteria2Group.GET("/config", hysteria2Handler.GetConfig)
        hysteria2Group.PUT("/config", hysteria2Handler.UpdateConfig)
        hysteria2Group.GET("/logs", hysteria2Handler.GetLogs)
        hysteria2Group.POST("/install", hysteria2Handler.Install)
        hysteria2Group.POST("/uninstall", hysteria2Handler.Uninstall)
        hysteria2Group.POST("/update", hysteria2Handler.Update)
        hysteria2Group.POST("/restart", hysteria2Handler.Restart)
        hysteria2Group.POST("/stop", hysteria2Handler.Stop)
        hysteria2Group.POST("/start", hysteria2Handler.Start)
        hysteria2Group.GET("/health", hysteria2Handler.CheckHealth)
        hysteria2Group.GET("/versions", hysteria2Handler.GetVersions)
        hysteria2Group.POST("/versions/install", hysteria2Handler.InstallVersion)
        hysteria2Group.GET("/config/backups", hysteria2Handler.GetConfigBackups)
        hysteria2Group.POST("/config/restore", hysteria2Handler.RestoreConfig)
    }
    
    // 配置管理API
    configHandler := v1.NewConfigHandler(cfg)
    configGroup := r.Group("/api/v1/config")
    {
        configGroup.GET("/whitelist", configHandler.GetWhitelist)
        configGroup.PUT("/whitelist", configHandler.UpdateWhitelist)
        configGroup.GET("/blacklist", configHandler.GetBlacklist)
        configGroup.PUT("/blacklist", configHandler.UpdateBlacklist)
    }
    
    r.Run(":8080")
}

func authMiddleware(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取客户端IP
        clientIP := c.ClientIP()
        
        // 检查黑名单
        for _, ip := range cfg.IPBlacklist {
            if ip == clientIP {
                c.JSON(403, gin.H{"error": "IP is blacklisted"})
                c.Abort()
                return
            }
        }
        
        // 如果设置了白名单，检查IP是否在白名单中
        if len(cfg.IPWhitelist) > 0 {
            allowed := false
            for _, ip := range cfg.IPWhitelist {
                if ip == clientIP {
                    allowed = true
                    break
                }
            }
            if !allowed {
                c.JSON(403, gin.H{"error": "IP not in whitelist"})
                c.Abort()
                return
            }
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