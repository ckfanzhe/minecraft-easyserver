package handlers

import (
	"net/http"
	"strconv"

	"minecraft-easyserver/services"
	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	logService *services.LogService
}

func NewLogHandler() *LogHandler {
	return &LogHandler{
		logService: services.NewLogService(),
	}
}

// GetLogs handles GET /api/logs
func (h *LogHandler) GetLogs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	logs := h.logService.GetLogs(limit)
	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"count": len(logs),
	})
}

// ClearLogs handles DELETE /api/logs
func (h *LogHandler) ClearLogs(c *gin.Context) {
	h.logService.ClearLogs()
	c.JSON(http.StatusOK, gin.H{
		"message": "Logs cleared successfully",
	})
}

// HandleWebSocket handles WebSocket connections for real-time logs
func (h *LogHandler) HandleWebSocket(c *gin.Context) {
	h.logService.HandleWebSocket(c.Writer, c.Request)
}