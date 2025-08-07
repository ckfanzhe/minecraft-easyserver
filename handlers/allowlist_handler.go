package handlers

import (
	"strings"

	"minecraft-easyserver/models"
	"minecraft-easyserver/services"

	"github.com/gin-gonic/gin"
)

// AllowlistHandler allowlist handler
type AllowlistHandler struct {
	allowlistService *services.AllowlistService
}

// NewAllowlistHandler creates a new allowlist handler
func NewAllowlistHandler() *AllowlistHandler {
	return &AllowlistHandler{
		allowlistService: services.NewAllowlistService(),
	}
}

// GetAllowlist gets allowlist
func (h *AllowlistHandler) GetAllowlist(c *gin.Context) {
	allowlist, err := h.allowlistService.GetAllowlist()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read allowlist: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"allowlist": allowlist})
}

// AddToAllowlist adds to allowlist
func (h *AllowlistHandler) AddToAllowlist(c *gin.Context) {
	var request struct {
		Name               string `json:"name"`
		IgnoresPlayerLimit bool   `json:"ignoresPlayerLimit"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	entry := models.AllowlistEntry{
		Name:               request.Name,
		IgnoresPlayerLimit: request.IgnoresPlayerLimit,
	}

	if err := h.allowlistService.AddToAllowlist(entry); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Added to allowlist: " + request.Name})
}

// RemoveFromAllowlist removes from allowlist
func (h *AllowlistHandler) RemoveFromAllowlist(c *gin.Context) {
	name := c.Param("name")

	if err := h.allowlistService.RemoveFromAllowlist(name); err != nil {
		if strings.Contains(err.Error(), "not in allowlist") {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Failed to save allowlist: " + err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Removed from allowlist: " + name})
}