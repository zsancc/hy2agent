package v1

import (
    "github.com/gin-gonic/gin"
    "agentapi/internal/service"
    "net/http"
    "strconv"
)

type Hysteria2Handler struct {
    hy2Service *service.Hysteria2Service
}

func NewHysteria2Handler() *Hysteria2Handler {
    return &Hysteria2Handler{
        hy2Service: service.NewHysteria2Service(),
    }
}

// 获取Hysteria2状态
func (h *Hysteria2Handler) GetStatus(c *gin.Context) {
    status, err := h.hy2Service.GetStatus()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, status)
}

// 安装Hysteria2
func (h *Hysteria2Handler) Install(c *gin.Context) {
    output, err := h.hy2Service.Install()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "output": output,
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "Hysteria2 installed successfully",
        "output": output,
    })
}

// 卸载Hysteria2
func (h *Hysteria2Handler) Uninstall(c *gin.Context) {
    output, err := h.hy2Service.Uninstall()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "output": output,
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "Hysteria2 uninstalled successfully",
        "output": output,
    })
}

// 更新Hysteria2
func (h *Hysteria2Handler) Update(c *gin.Context) {
    output, err := h.hy2Service.Update()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "output": output,
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "Hysteria2 updated successfully",
        "output": output,
    })
}

// 获取配置
func (h *Hysteria2Handler) GetConfig(c *gin.Context) {
    config, err := h.hy2Service.GetConfig()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"config": config})
}

// 更新配置
func (h *Hysteria2Handler) UpdateConfig(c *gin.Context) {
    var req struct {
        Config string `json:"config" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.hy2Service.UpdateConfig(req.Config); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Config updated successfully"})
}

// 获取日志
func (h *Hysteria2Handler) GetLogs(c *gin.Context) {
    var opts service.LogOptions
    // 从查询参数获取选项
    if lines := c.Query("lines"); lines != "" {
        if n, err := strconv.Atoi(lines); err == nil {
            opts.Lines = n
        }
    }
    opts.Since = c.Query("since")  // 如 "5m", "2h"
    opts.Level = c.Query("level")  // 如 "info", "error"

    logs, err := h.hy2Service.GetLogs(&opts)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// 启动服务
func (h *Hysteria2Handler) Start(c *gin.Context) {
    if err := h.hy2Service.Start(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Hysteria2 service started"})
}

// 停止服务
func (h *Hysteria2Handler) Stop(c *gin.Context) {
    if err := h.hy2Service.Stop(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Hysteria2 service stopped"})
}

// 重启服务
func (h *Hysteria2Handler) Restart(c *gin.Context) {
    if err := h.hy2Service.Restart(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Hysteria2 service restarted"})
}

// 健康检查
func (h *Hysteria2Handler) CheckHealth(c *gin.Context) {
    health, err := h.hy2Service.CheckHealth()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, health)
}

// 获取可用版本
func (h *Hysteria2Handler) GetVersions(c *gin.Context) {
    versions, err := h.hy2Service.GetAvailableVersions()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"versions": versions})
}

// 安装指定版本
func (h *Hysteria2Handler) InstallVersion(c *gin.Context) {
    var req struct {
        Version string `json:"version" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.hy2Service.InstallVersion(req.Version); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Version installed successfully"})
}

// 获取配置备份列表
func (h *Hysteria2Handler) GetConfigBackups(c *gin.Context) {
    backups, err := h.hy2Service.GetConfigBackups()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"backups": backups})
}

// 恢复配置备份
func (h *Hysteria2Handler) RestoreConfig(c *gin.Context) {
    var req struct {
        Backup string `json:"backup" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.hy2Service.RestoreConfig(req.Backup); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Config restored successfully"})
}