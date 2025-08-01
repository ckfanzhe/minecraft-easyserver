package handlers

import (
	"fmt"
	"strings"

	"minecraft-easyserver/services"

	"github.com/gin-gonic/gin"
)

// PermissionHandler permission handler
type PermissionHandler struct {
	permissionService *services.PermissionService
}

// NewPermissionHandler creates a new permission handler
func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{
		permissionService: services.NewPermissionService(),
	}
}

// GetPermissions gets permissions
func (h *PermissionHandler) GetPermissions(c *gin.Context) {
	permissions, err := h.permissionService.GetPermissions()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read permissions: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"permissions": permissions})
}

// UpdatePermission updates permission
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	var request struct {
		Name  string `json:"name"`
		Level string `json:"level"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.permissionService.UpdatePermission(request.Name, request.Level); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Set %s permission to %s", request.Name, request.Level)})
}

// RemovePermission removes permission
func (h *PermissionHandler) RemovePermission(c *gin.Context) {
	name := c.Param("name")

	if err := h.permissionService.RemovePermission(name); err != nil {
		if strings.Contains(err.Error(), "permission not found") {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Failed to save permissions: " + err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Permission removed: " + name})
}