package handlers

import (
	"minecraft-easyserver/services"

	"github.com/gin-gonic/gin"
)

// ServerHandler server handler
type ServerHandler struct {
	serverService *services.ServerService
}

// NewServerHandler creates a new server handler
func NewServerHandler() *ServerHandler {
	return &ServerHandler{
		serverService: services.NewServerService(),
	}
}

// GetStatus gets server status
func (h *ServerHandler) GetStatus(c *gin.Context) {
	status := h.serverService.GetStatus()
	c.JSON(200, status)
}

// StartServer starts the server
func (h *ServerHandler) StartServer(c *gin.Context) {
	if err := h.serverService.Start(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Server started successfully"})
}

// StopServer stops the server
func (h *ServerHandler) StopServer(c *gin.Context) {
	if err := h.serverService.Stop(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Server stopped"})
}

// RestartServer restarts the server
func (h *ServerHandler) RestartServer(c *gin.Context) {
	if err := h.serverService.Restart(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Server restarted successfully"})
}