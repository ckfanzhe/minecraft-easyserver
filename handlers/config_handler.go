package handlers

import (
	"minecraft-easyserver/models"
	"minecraft-easyserver/services"

	"github.com/gin-gonic/gin"
)

// ConfigHandler configuration handler
type ConfigHandler struct {
	configService *services.ConfigService
}

// NewConfigHandler creates a new configuration handler
func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		configService: services.NewConfigService(),
	}
}

// GetConfig gets server configuration
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	config, err := h.configService.GetConfig()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read configuration: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"config": config})
}

// UpdateConfig updates server configuration
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	var request struct {
		Config models.ServerConfig `json:"config"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.configService.UpdateConfig(request.Config); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save configuration: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Configuration saved, restart server to take effect"})
}