package v1

import (
	"hy2agent/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	statusService *service.StatusService
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{
		statusService: service.NewStatusService(),
	}
}

// 获取内存信息
func (h *SystemHandler) GetMemory(c *gin.Context) {
	memInfo, err := h.statusService.GetMemoryInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memInfo)
}

// 获取磁盘信息
func (h *SystemHandler) GetDisk(c *gin.Context) {
	diskInfo, err := h.statusService.GetDiskInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, diskInfo)
}

// 获取网络信息
func (h *SystemHandler) GetNetwork(c *gin.Context) {
	netInfo, err := h.statusService.GetNetworkInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, netInfo)
}

// 获取系统信息
func (h *SystemHandler) GetInfo(c *gin.Context) {
	sysInfo, err := h.statusService.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sysInfo)
}
