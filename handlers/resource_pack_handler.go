package handlers

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"minecraft-easyserver/services"

	"github.com/gin-gonic/gin"
)

// ResourcePackHandler resource pack handler
type ResourcePackHandler struct {
	resourcePackService *services.ResourcePackService
}

// NewResourcePackHandler creates a new resource pack handler
func NewResourcePackHandler() *ResourcePackHandler {
	return &ResourcePackHandler{
		resourcePackService: services.NewResourcePackService(),
	}
}

// GetResourcePacks gets resource pack list
func (h *ResourcePackHandler) GetResourcePacks(c *gin.Context) {
	packs, err := h.resourcePackService.GetResourcePacks()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read resource pack list: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"resource_packs": packs})
}

// UploadResourcePack uploads resource pack
func (h *ResourcePackHandler) UploadResourcePack(c *gin.Context) {
	file, header, err := c.Request.FormFile("resource_pack")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}
	defer file.Close()

	// Check file extension
	filename := header.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".zip") &&
		!strings.HasSuffix(strings.ToLower(filename), ".mcpack") {
		c.JSON(400, gin.H{"error": "Only .zip and .mcpack formats are supported"})
		return
	}

	// Create temporary directory
	tempDir := os.TempDir()
	tempPath := filepath.Join(tempDir, filename)

	// Save uploaded file to temporary location
	out, err := os.Create(tempPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file: " + err.Error()})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file: " + err.Error()})
		return
	}
	out.Close()

	// Upload and extract resource pack
	packInfo, err := h.resourcePackService.UploadResourcePack(tempPath, filename)
	if err != nil {
		// Clean up temporary file
		os.Remove(tempPath)
		c.JSON(500, gin.H{"error": "Failed to process resource pack: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":       "Resource pack uploaded successfully",
		"resource_pack": packInfo,
	})
}

// ActivateResourcePack activates resource pack
func (h *ResourcePackHandler) ActivateResourcePack(c *gin.Context) {
	packUUID := c.Param("uuid")

	if err := h.resourcePackService.ActivateResourcePack(packUUID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(404, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "already activated") {
			c.JSON(400, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "system resource pack") {
			c.JSON(403, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Failed to activate resource pack: " + err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Resource pack activated, restart server to take effect"})
}

// DeactivateResourcePack deactivates resource pack
func (h *ResourcePackHandler) DeactivateResourcePack(c *gin.Context) {
	packUUID := c.Param("uuid")

	if err := h.resourcePackService.DeactivateResourcePack(packUUID); err != nil {
		if strings.Contains(err.Error(), "not activated") {
			c.JSON(400, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "system resource pack") {
			c.JSON(403, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Failed to deactivate resource pack: " + err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Resource pack deactivated, restart server to take effect"})
}

// DeleteResourcePack deletes resource pack
func (h *ResourcePackHandler) DeleteResourcePack(c *gin.Context) {
	packUUID := c.Param("uuid")

	if err := h.resourcePackService.DeleteResourcePack(packUUID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(404, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "system resource pack") {
			c.JSON(403, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Failed to delete resource pack: " + err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Resource pack deleted successfully"})
}