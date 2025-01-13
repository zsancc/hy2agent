package v1

import (
    "github.com/gin-gonic/gin"
    "agentapi/internal/service"
    "net/http"
)

type StatusHandler struct {
    statusService *service.StatusService
}

func NewStatusHandler() *StatusHandler {
    return &StatusHandler{
        statusService: service.NewStatusService(),
    }
}

func (h *StatusHandler) GetStatus(c *gin.Context) {
    status, err := h.statusService.GetSystemStatus()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, status)
} 