package handlers

import (
	"net/http"
	"strconv"

	"minecraft-easyserver/models"
	"minecraft-easyserver/services"
	"github.com/gin-gonic/gin"
)

type InteractionHandler struct {
	interactionService *services.InteractionService
}

func NewInteractionHandler() *InteractionHandler {
	return &InteractionHandler{
		interactionService: services.GetInteractionService(),
	}
}

// GetStatus handles GET /api/interaction/status
func (h *InteractionHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"enabled": h.interactionService.IsEnabled(),
		"platform": "linux", // Could be dynamic based on runtime.GOOS
	})
}

// SendCommand handles POST /api/interaction/command
func (h *InteractionHandler) SendCommand(c *gin.Context) {
	var req models.ServerCommand
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if !h.interactionService.IsEnabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Server interaction is not supported on this platform",
		})
		return
	}

	// Validate command
	if err := h.interactionService.ValidateCommand(req.Command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Send command
	if err := h.interactionService.SendCommand(req.Command); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Command sent successfully",
		"command": req.Command,
	})
}

// GetCommandHistory handles GET /api/interaction/history
func (h *InteractionHandler) GetCommandHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	history := h.interactionService.GetCommandHistory(limit)
	c.JSON(http.StatusOK, gin.H{
		"history": history,
		"count": len(history),
	})
}

// ClearHistory handles DELETE /api/interaction/history
func (h *InteractionHandler) ClearHistory(c *gin.Context) {
	h.interactionService.ClearHistory()
	c.JSON(http.StatusOK, gin.H{
		"message": "Command history cleared successfully",
	})
}