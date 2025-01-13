package v1

import (
	"hy2agent/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
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
