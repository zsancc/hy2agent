package v1

import (
	"hy2agent/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	cfg *config.Config
}

func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{cfg: cfg}
}

// 获取IP白名单
func (h *ConfigHandler) GetWhitelist(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"whitelist": h.cfg.IPWhitelist})
}

// 更新IP白名单
func (h *ConfigHandler) UpdateWhitelist(c *gin.Context) {
	var req struct {
		IPs []string `json:"ips" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.cfg.IPWhitelist = req.IPs
	if err := config.SaveConfig(h.cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Whitelist updated"})
}
