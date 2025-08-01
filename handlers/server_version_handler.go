package handlers

import (
	"minecraft-easyserver/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServerVersionHandler handles server version management requests
type ServerVersionHandler struct {
	service *services.ServerVersionService
}

// NewServerVersionHandler creates a new server version handler
func NewServerVersionHandler() *ServerVersionHandler {
	return &ServerVersionHandler{
		service: services.NewServerVersionService(),
	}
}

// GetVersions returns available server versions
func (h *ServerVersionHandler) GetVersions(c *gin.Context) {
	versions := h.service.GetAvailableVersions()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    versions,
	})
}

// DownloadVersion downloads a specific server version
func (h *ServerVersionHandler) DownloadVersion(c *gin.Context) {
	version := c.Param("version")
	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Version parameter is required",
		})
		return
	}

	err := h.service.DownloadVersion(version)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Download started",
	})
}

// GetDownloadProgress returns download progress for a version
func (h *ServerVersionHandler) GetDownloadProgress(c *gin.Context) {
	version := c.Param("version")
	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Version parameter is required",
		})
		return
	}

	progress, exists := h.service.GetDownloadProgress(version)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No download progress found for this version",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    progress,
	})
}

// ActivateVersion activates a specific server version
func (h *ServerVersionHandler) ActivateVersion(c *gin.Context) {
	version := c.Param("version")
	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Version parameter is required",
		})
		return
	}

	err := h.service.ActivateVersion(version)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Version activated successfully",
	})
}

// UpdateVersionConfig updates the server version configuration from GitHub
func (h *ServerVersionHandler) UpdateVersionConfig(c *gin.Context) {
	err := h.service.UpdateVersionConfigFromGitHub()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Return updated versions
	versions := h.service.GetAvailableVersions()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Version configuration updated successfully",
		"data":    versions,
	})
}