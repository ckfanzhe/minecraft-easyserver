package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"minecraft-easyserver/services"
	"minecraft-easyserver/utils"

	"github.com/gin-gonic/gin"
)

// WorldHandler world handler
type WorldHandler struct {
	worldService *services.WorldService
}

// NewWorldHandler creates a new world handler
func NewWorldHandler() *WorldHandler {
	return &WorldHandler{
		worldService: services.NewWorldService(),
	}
}

// GetWorlds gets world list
func (h *WorldHandler) GetWorlds(c *gin.Context) {
	worlds, err := h.worldService.GetWorlds()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read world list: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"worlds": worlds})
}



// UploadWorld uploads world
func (h *WorldHandler) UploadWorld(c *gin.Context) {
	file, header, err := c.Request.FormFile("world")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}
	defer file.Close()

	// Check file extension
	filename := header.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".zip") &&
		!strings.HasSuffix(strings.ToLower(filename), ".mcworld") {
		c.JSON(400, gin.H{"error": "Only .zip and .mcworld formats are supported"})
		return
	}

	// Get bedrock path
	wd, err := os.Getwd()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get working directory: " + err.Error()})
		return
	}
	bedrockPath := filepath.Join(wd, "bedrock-server", "bedrock-server-1.21.95.1")

	// Save uploaded file
	worldsPath := filepath.Join(bedrockPath, "worlds")
	os.MkdirAll(worldsPath, 0755)

	uploadPath := filepath.Join(worldsPath, filename)
	out, err := os.Create(uploadPath)
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

	// Close file handle for subsequent operations
	out.Close()

	// Extract file
	var extractedWorldName string
	if strings.HasSuffix(strings.ToLower(filename), ".zip") || strings.HasSuffix(strings.ToLower(filename), ".mcworld") {
		// Get filename without extension as world name
		extractedWorldName = strings.TrimSuffix(filename, filepath.Ext(filename))
		extractPath := filepath.Join(worldsPath, extractedWorldName)

		// Extract file
		if err := utils.ExtractZip(uploadPath, extractPath); err != nil {
			// If extraction fails, delete uploaded file
			os.Remove(uploadPath)
			c.JSON(500, gin.H{"error": "Failed to extract file: " + err.Error()})
			return
		}

		// Delete original compressed file after successful extraction
		if err := os.Remove(uploadPath); err != nil {
			// Log warning but don't affect main flow
			fmt.Printf("Warning: Failed to delete compressed file: %v\n", err)
		}

		c.JSON(200, gin.H{"message": fmt.Sprintf("World file uploaded and extracted successfully: %s", extractedWorldName)})
	} else {
		c.JSON(200, gin.H{"message": "World file uploaded successfully: " + filename})
	}
}

// DeleteWorld deletes world
func (h *WorldHandler) DeleteWorld(c *gin.Context) {
	worldName := c.Param("name")

	if err := h.worldService.DeleteWorld(worldName); err != nil {
		// Return different status codes based on error type
		if strings.Contains(err.Error(), "world not found") {
			c.JSON(404, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "world name cannot be empty") {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Failed to delete world: " + err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "World deleted: " + worldName + ", configuration file updated"})
}

// ActivateWorld activates world
func (h *WorldHandler) ActivateWorld(c *gin.Context) {
	worldName := c.Param("name")

	if err := h.worldService.ActivateWorld(worldName); err != nil {
		c.JSON(500, gin.H{"error": "Failed to activate world: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "World activated: " + worldName + ", restart server to take effect"})
}