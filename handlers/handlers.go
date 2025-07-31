package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"bedrock-easyserver/models"
	"bedrock-easyserver/services"
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
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.allowlistService.AddToAllowlist(request.Name); err != nil {
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

// extractZip extracts zip file
func extractZip(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Ensure destination directory exists
	os.MkdirAll(dest, 0755)

	// Extract files
	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)

		// Check path security to prevent directory traversal attacks
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.FileInfo().Mode())
			continue
		}

		// Create file directory
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// Create file
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
	}

	return nil
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
		if err := extractZip(uploadPath, extractPath); err != nil {
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
		"message":      "Resource pack uploaded successfully",
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
		} else {
			c.JSON(500, gin.H{"error": "Failed to delete resource pack: " + err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Resource pack deleted successfully"})
}