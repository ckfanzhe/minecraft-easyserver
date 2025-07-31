package services

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"bedrock-easyserver/models"
)

var (
	serverProcess *exec.Cmd
	serverMutex   sync.Mutex
	bedrockPath   string
)

// InitBedrockPath initializes bedrock path
func InitBedrockPath(path string) error {
	if path == "" {
		return fmt.Errorf("bedrock path cannot be empty")
	}
	
	// If it's a relative path, convert to absolute path
	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		path = filepath.Join(wd, path)
	}
	
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("bedrock path does not exist: %s", path)
	}
	
	bedrockPath = path
	return nil
}

// SetBedrockPath sets bedrock path (mainly for testing)
func SetBedrockPath(path string) {
	bedrockPath = path
}

// ServerService server service
type ServerService struct{}

// NewServerService creates a new server service instance
func NewServerService() *ServerService {
	return &ServerService{}
}

// GetStatus gets server status
func (s *ServerService) GetStatus() models.ServerStatus {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverProcess == nil || serverProcess.Process == nil {
		return models.ServerStatus{
			Status:  "stopped",
			Message: "Server not running",
		}
	}

	// Check if process is still running
	process, err := os.FindProcess(serverProcess.Process.Pid)
	if err != nil {
		serverProcess = nil
		return models.ServerStatus{
			Status:  "stopped",
			Message: "Server not running",
		}
	}

	// On Windows, simply check if process exists
	// If process has ended, FindProcess will still return a Process object
	// We can try sending signal 0 to check if process is really running
	if process != nil {
		return models.ServerStatus{
			Status:  "running",
			Message: "Server is running",
			PID:     serverProcess.Process.Pid,
		}
	}

	serverProcess = nil
	return models.ServerStatus{
		Status:  "stopped",
		Message: "Server not running",
	}
}

// Start starts server
func (s *ServerService) Start() error {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverProcess != nil && serverProcess.Process != nil {
		return fmt.Errorf("server is already running")
	}

	exePath := filepath.Join(bedrockPath, "bedrock_server.exe")
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		return fmt.Errorf("bedrock_server.exe file not found")
	}

	serverProcess = exec.Command(exePath)
	serverProcess.Dir = bedrockPath

	if err := serverProcess.Start(); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}

// Stop stops server
func (s *ServerService) Stop() error {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverProcess == nil || serverProcess.Process == nil {
		return fmt.Errorf("server not running")
	}

	if err := serverProcess.Process.Kill(); err != nil {
		return fmt.Errorf("failed to stop server: %v", err)
	}

	serverProcess.Wait()
	serverProcess = nil
	return nil
}

// Restart restarts server
func (s *ServerService) Restart() error {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	// Stop first
	if serverProcess != nil && serverProcess.Process != nil {
		serverProcess.Process.Kill()
		serverProcess.Wait()
		serverProcess = nil
	}

	// Wait one second
	time.Sleep(time.Second)

	// Restart
	exePath := filepath.Join(bedrockPath, "bedrock_server.exe")
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		return fmt.Errorf("bedrock_server.exe file not found")
	}

	serverProcess = exec.Command(exePath)
	serverProcess.Dir = bedrockPath

	if err := serverProcess.Start(); err != nil {
		return fmt.Errorf("failed to restart server: %v", err)
	}

	return nil
}

// ConfigService configuration service
type ConfigService struct{}

// NewConfigService creates a new configuration service instance
func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// GetConfig gets server configuration
func (c *ConfigService) GetConfig() (models.ServerConfig, error) {
	configPath := filepath.Join(bedrockPath, "server.properties")
	return readServerProperties(configPath)
}

// UpdateConfig updates server configuration
func (c *ConfigService) UpdateConfig(config models.ServerConfig) error {
	configPath := filepath.Join(bedrockPath, "server.properties")
	return writeServerProperties(configPath, config)
}

// Read server.properties
func readServerProperties(path string) (models.ServerConfig, error) {
	config := models.ServerConfig{}

	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "server-name":
			config.ServerName = value
		case "gamemode":
			config.Gamemode = value
		case "difficulty":
			config.Difficulty = value
		case "max-players":
			if i, err := strconv.Atoi(value); err == nil {
				config.MaxPlayers = i
			}
		case "server-port":
			if i, err := strconv.Atoi(value); err == nil {
				config.ServerPort = i
			}
		case "allow-cheats":
			config.AllowCheats = value == "true"
		case "allow-list":
			config.AllowList = value == "true"
		case "online-mode":
			config.OnlineMode = value == "true"
		case "level-name":
			config.LevelName = value
		case "default-player-permission-level":
			config.DefaultPlayerPermission = value
		}
	}

	return config, scanner.Err()
}

// Write server.properties
func writeServerProperties(path string, config models.ServerConfig) error {
	// Read original file to preserve comments
	originalLines := []string{}
	if file, err := os.Open(path); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			originalLines = append(originalLines, scanner.Text())
		}
		file.Close()
	}

	// Create new file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write configuration
	configMap := map[string]string{
		"server-name":                     config.ServerName,
		"gamemode":                        config.Gamemode,
		"difficulty":                      config.Difficulty,
		"max-players":                     strconv.Itoa(config.MaxPlayers),
		"server-port":                     strconv.Itoa(config.ServerPort),
		"allow-cheats":                    strconv.FormatBool(config.AllowCheats),
		"allow-list":                      strconv.FormatBool(config.AllowList),
		"online-mode":                     strconv.FormatBool(config.OnlineMode),
		"level-name":                      config.LevelName,
		"default-player-permission-level": config.DefaultPlayerPermission,
	}

	// Process each line
	for _, line := range originalLines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			writer.WriteString(line + "\n")
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			writer.WriteString(line + "\n")
			continue
		}

		key := strings.TrimSpace(parts[0])
		if newValue, exists := configMap[key]; exists {
			writer.WriteString(key + "=" + newValue + "\n")
			delete(configMap, key)
		} else {
			writer.WriteString(line + "\n")
		}
	}

	// Add any new configuration items
	for key, value := range configMap {
		writer.WriteString(key + "=" + value + "\n")
	}

	return nil
}

// AllowlistService allowlist service
type AllowlistService struct{}

// NewAllowlistService creates a new allowlist service instance
func NewAllowlistService() *AllowlistService {
	return &AllowlistService{}
}

// GetAllowlist gets allowlist
func (a *AllowlistService) GetAllowlist() ([]string, error) {
	allowlistPath := filepath.Join(bedrockPath, "allowlist.json")
	allowlist, err := readAllowlist(allowlistPath)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range allowlist {
		names = append(names, entry.Name)
	}

	return names, nil
}

// AddToAllowlist adds to allowlist
func (a *AllowlistService) AddToAllowlist(name string) error {
	if name == "" {
		return fmt.Errorf("player name cannot be empty")
	}

	allowlistPath := filepath.Join(bedrockPath, "allowlist.json")
	allowlist, err := readAllowlist(allowlistPath)
	if err != nil {
		allowlist = []models.AllowlistEntry{}
	}

	// Check if already exists
	for _, entry := range allowlist {
		if entry.Name == name {
			return fmt.Errorf("player already in allowlist")
		}
	}

	// Add new entry
	newEntry := models.AllowlistEntry{
		Name:               name,
		IgnoresPlayerLimit: false,
	}
	allowlist = append(allowlist, newEntry)

	return writeAllowlist(allowlistPath, allowlist)
}

// RemoveFromAllowlist removes from allowlist
func (a *AllowlistService) RemoveFromAllowlist(name string) error {
	if name == "" {
		return fmt.Errorf("player name cannot be empty")
	}

	allowlistPath := filepath.Join(bedrockPath, "allowlist.json")
	allowlist, err := readAllowlist(allowlistPath)
	if err != nil {
		return err
	}

	// Remove entry
	var newAllowlist []models.AllowlistEntry
	found := false
	for _, entry := range allowlist {
		if entry.Name != name {
			newAllowlist = append(newAllowlist, entry)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("player not in allowlist")
	}

	return writeAllowlist(allowlistPath, newAllowlist)
}

// Read allowlist.json
func readAllowlist(path string) ([]models.AllowlistEntry, error) {
	var allowlist []models.AllowlistEntry

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return allowlist, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &allowlist); err != nil {
		return nil, err
	}

	return allowlist, nil
}

// Write allowlist.json
func writeAllowlist(path string, allowlist []models.AllowlistEntry) error {
	data, err := json.MarshalIndent(allowlist, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// PermissionService permission service
type PermissionService struct{}

// NewPermissionService creates a new permission service instance
func NewPermissionService() *PermissionService {
	return &PermissionService{}
}

// GetPermissions gets permissions
func (p *PermissionService) GetPermissions() ([]map[string]interface{}, error) {
	permissionsPath := filepath.Join(bedrockPath, "permissions.json")
	return readPermissions(permissionsPath)
}

// UpdatePermission updates permission
func (p *PermissionService) UpdatePermission(name, level string) error {
	if name == "" {
		return fmt.Errorf("player name cannot be empty")
	}

	validLevels := map[string]bool{
		"visitor":  true,
		"member":   true,
		"operator": true,
	}

	if !validLevels[level] {
		return fmt.Errorf("invalid permission level")
	}

	permissionsPath := filepath.Join(bedrockPath, "permissions.json")
	permissions, err := readPermissions(permissionsPath)
	if err != nil {
		permissions = []map[string]interface{}{}
	}

	// Find and update or add permission
	found := false
	for i, perm := range permissions {
		if playerName, ok := perm["name"].(string); ok && playerName == name {
			permissions[i]["level"] = level
			found = true
			break
		}
	}

	if !found {
		newPerm := map[string]interface{}{
			"name":  name,
			"level": level,
		}
		permissions = append(permissions, newPerm)
	}

	return writePermissions(permissionsPath, permissions)
}

// RemovePermission removes permission
func (p *PermissionService) RemovePermission(name string) error {
	if name == "" {
		return fmt.Errorf("player name cannot be empty")
	}

	permissionsPath := filepath.Join(bedrockPath, "permissions.json")
	permissions, err := readPermissions(permissionsPath)
	if err != nil {
		return err
	}

	// Remove permission
	var newPermissions []map[string]interface{}
	found := false
	for _, perm := range permissions {
		if playerName, ok := perm["name"].(string); ok && playerName != name {
			newPermissions = append(newPermissions, perm)
		} else if ok {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("player permission not found")
	}

	return writePermissions(permissionsPath, newPermissions)
}

// Read permissions.json
func readPermissions(path string) ([]map[string]interface{}, error) {
	var permissions []map[string]interface{}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return permissions, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

// Write permissions.json
func writePermissions(path string, permissions []map[string]interface{}) error {
	data, err := json.MarshalIndent(permissions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// WorldService world service
type WorldService struct{}

// NewWorldService creates a new world service instance
func NewWorldService() *WorldService {
	return &WorldService{}
}

// GetWorlds gets world list
func (w *WorldService) GetWorlds() ([]models.WorldInfo, error) {
	worldsPath := filepath.Join(bedrockPath, "worlds")
	return getWorldsList(worldsPath)
}

// DeleteWorld deletes world
func (w *WorldService) DeleteWorld(worldName string) error {
	if worldName == "" {
		return fmt.Errorf("world name cannot be empty")
	}

	// Check if world exists
	worldPath := filepath.Join(bedrockPath, "worlds", worldName)
	if _, err := os.Stat(worldPath); os.IsNotExist(err) {
		return fmt.Errorf("world not found: %s", worldName)
	}

	// Check if it's the currently active world
	configPath := filepath.Join(bedrockPath, "server.properties")
	config, err := readServerProperties(configPath)
	isActiveWorld := false
	if err == nil && config.LevelName == worldName {
		isActiveWorld = true
	}

	// Delete world folder
	if err := os.RemoveAll(worldPath); err != nil {
		return fmt.Errorf("failed to delete world files: %v", err)
	}

	// If deleting the currently active world, need to update configuration file
	if isActiveWorld {
		// Get remaining world list
		worldsPath := filepath.Join(bedrockPath, "worlds")
		remainingWorlds, err := getWorldsList(worldsPath)
		if err != nil {
			return fmt.Errorf("failed to get remaining world list: %v", err)
		}

		// If there are other worlds, activate the first one; otherwise set to default world name
		if len(remainingWorlds) > 0 {
			config.LevelName = remainingWorlds[0].Name
		} else {
			// When no other worlds exist, set to default world name
			config.LevelName = "Bedrock level"
		}

		// Update configuration file
		if err := writeServerProperties(configPath, config); err != nil {
			return fmt.Errorf("failed to update configuration file: %v", err)
		}
	}

	return nil
}

// ActivateWorld activates world
func (w *WorldService) ActivateWorld(worldName string) error {
	if worldName == "" {
		return fmt.Errorf("world name cannot be empty")
	}

	// Update level-name in server.properties
	configPath := filepath.Join(bedrockPath, "server.properties")
	config, err := readServerProperties(configPath)
	if err != nil {
		return err
	}

	config.LevelName = worldName
	return writeServerProperties(configPath, config)
}

// Get world list
func getWorldsList(worldsPath string) ([]models.WorldInfo, error) {
	var worlds []models.WorldInfo

	// Read currently active world
	configPath := filepath.Join(bedrockPath, "server.properties")
	config, err := readServerProperties(configPath)
	activeWorld := ""
	if err == nil {
		activeWorld = config.LevelName
	}

	entries, err := os.ReadDir(worldsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return worlds, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			world := models.WorldInfo{
				Name:   entry.Name(),
				Active: entry.Name() == activeWorld,
			}
			worlds = append(worlds, world)
		}
	}

	return worlds, nil
}

// ResourcePackService resource pack service
type ResourcePackService struct{}

// NewResourcePackService creates a new resource pack service instance
func NewResourcePackService() *ResourcePackService {
	return &ResourcePackService{}
}

// UploadResourcePack uploads and extracts resource pack
func (r *ResourcePackService) UploadResourcePack(zipPath, fileName string) (*models.ResourcePackInfo, error) {
	// Create resource_packs directory if it doesn't exist
	resourcePacksPath := filepath.Join(bedrockPath, "resource_packs")
	if err := os.MkdirAll(resourcePacksPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create resource_packs directory: %v", err)
	}

	// Extract zip file
	extractPath := filepath.Join(resourcePacksPath, strings.TrimSuffix(fileName, filepath.Ext(fileName)))
	if err := extractZip(zipPath, extractPath); err != nil {
		return nil, fmt.Errorf("failed to extract resource pack: %v", err)
	}

	// Read manifest.json
	manifestPath := filepath.Join(extractPath, "manifest.json")
	manifest, err := readResourcePackManifest(manifestPath)
	if err != nil {
		// Clean up extracted files if manifest reading fails
		os.RemoveAll(extractPath)
		return nil, fmt.Errorf("failed to read manifest.json: %v", err)
	}

	// Delete original zip file
	os.Remove(zipPath)

	// Return resource pack information
	packInfo := &models.ResourcePackInfo{
		Name:        manifest.Header.Name,
		UUID:        manifest.Header.UUID,
		Version:     manifest.Header.Version,
		Description: manifest.Header.Description,
		FolderName:  filepath.Base(extractPath),
		Active:      false, // Initially not active
	}

	return packInfo, nil
}

// GetResourcePacks gets resource pack list
func (r *ResourcePackService) GetResourcePacks() ([]models.ResourcePackInfo, error) {
	resourcePacksPath := filepath.Join(bedrockPath, "resource_packs")
	var packs []models.ResourcePackInfo

	// Get currently active world
	configPath := filepath.Join(bedrockPath, "server.properties")
	config, err := readServerProperties(configPath)
	if err != nil {
		return nil, err
	}

	// Read active resource packs from world configuration
	activePackUUIDs := make(map[string]bool)
	if config.LevelName != "" {
		worldResourcePacksPath := filepath.Join(bedrockPath, "worlds", config.LevelName, "world_resource_packs.json")
		activePacks, err := readWorldResourcePacks(worldResourcePacksPath)
		if err == nil {
			for _, pack := range activePacks {
				activePackUUIDs[pack.PackID] = true
			}
		}
	}

	// Read all resource pack directories
	entries, err := os.ReadDir(resourcePacksPath)
	if err != nil {
		if os.IsNotExist(err) {
			return packs, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			manifestPath := filepath.Join(resourcePacksPath, entry.Name(), "manifest.json")
			manifest, err := readResourcePackManifest(manifestPath)
			if err != nil {
				continue // Skip invalid resource packs
			}

			packInfo := models.ResourcePackInfo{
				Name:        manifest.Header.Name,
				UUID:        manifest.Header.UUID,
				Version:     manifest.Header.Version,
				Description: manifest.Header.Description,
				FolderName:  entry.Name(),
				Active:      activePackUUIDs[manifest.Header.UUID],
			}
			packs = append(packs, packInfo)
		}
	}

	return packs, nil
}

// ActivateResourcePack activates resource pack for current world
func (r *ResourcePackService) ActivateResourcePack(packUUID string) error {
	// Get current world
	configPath := filepath.Join(bedrockPath, "server.properties")
	config, err := readServerProperties(configPath)
	if err != nil {
		return err
	}

	if config.LevelName == "" {
		return fmt.Errorf("no active world found")
	}

	// Find resource pack by UUID
	resourcePacksPath := filepath.Join(bedrockPath, "resource_packs")
	entries, err := os.ReadDir(resourcePacksPath)
	if err != nil {
		return fmt.Errorf("failed to read resource packs directory: %v", err)
	}

	var targetPack *models.ResourcePackManifest
	for _, entry := range entries {
		if entry.IsDir() {
			manifestPath := filepath.Join(resourcePacksPath, entry.Name(), "manifest.json")
			manifest, err := readResourcePackManifest(manifestPath)
			if err != nil {
				continue
			}
			if manifest.Header.UUID == packUUID {
				targetPack = &manifest
				break
			}
		}
	}

	if targetPack == nil {
		return fmt.Errorf("resource pack not found: %s", packUUID)
	}

	// Read current world resource packs configuration
	worldResourcePacksPath := filepath.Join(bedrockPath, "worlds", config.LevelName, "world_resource_packs.json")
	activePacks, err := readWorldResourcePacks(worldResourcePacksPath)
	if err != nil {
		activePacks = []models.WorldResourcePack{}
	}

	// Check if already activated
	for _, pack := range activePacks {
		if pack.PackID == packUUID {
			return fmt.Errorf("resource pack already activated")
		}
	}

	// Add new resource pack
	newPack := models.WorldResourcePack{
		PackID:  targetPack.Header.UUID,
		Version: targetPack.Header.Version,
	}
	activePacks = append(activePacks, newPack)

	// Write back to file
	return writeWorldResourcePacks(worldResourcePacksPath, activePacks)
}

// DeactivateResourcePack deactivates resource pack for current world
func (r *ResourcePackService) DeactivateResourcePack(packUUID string) error {
	// Get current world
	configPath := filepath.Join(bedrockPath, "server.properties")
	config, err := readServerProperties(configPath)
	if err != nil {
		return err
	}

	if config.LevelName == "" {
		return fmt.Errorf("no active world found")
	}

	// Read current world resource packs configuration
	worldResourcePacksPath := filepath.Join(bedrockPath, "worlds", config.LevelName, "world_resource_packs.json")
	activePacks, err := readWorldResourcePacks(worldResourcePacksPath)
	if err != nil {
		return fmt.Errorf("failed to read world resource packs configuration: %v", err)
	}

	// Remove resource pack
	var newActivePacks []models.WorldResourcePack
	found := false
	for _, pack := range activePacks {
		if pack.PackID != packUUID {
			newActivePacks = append(newActivePacks, pack)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("resource pack not activated")
	}

	// Write back to file
	return writeWorldResourcePacks(worldResourcePacksPath, newActivePacks)
}

// DeleteResourcePack deletes resource pack
func (r *ResourcePackService) DeleteResourcePack(packUUID string) error {
	// First deactivate if active
	r.DeactivateResourcePack(packUUID) // Ignore error as pack might not be active

	// Find and delete resource pack directory
	resourcePacksPath := filepath.Join(bedrockPath, "resource_packs")
	entries, err := os.ReadDir(resourcePacksPath)
	if err != nil {
		return fmt.Errorf("failed to read resource packs directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			manifestPath := filepath.Join(resourcePacksPath, entry.Name(), "manifest.json")
			manifest, err := readResourcePackManifest(manifestPath)
			if err != nil {
				continue
			}
			if manifest.Header.UUID == packUUID {
				packPath := filepath.Join(resourcePacksPath, entry.Name())
				return os.RemoveAll(packPath)
			}
		}
	}

	return fmt.Errorf("resource pack not found: %s", packUUID)
}

// Helper functions

// extractZip extracts zip file to target directory
func extractZip(src, dest string) error {
	// Open zip file for reading
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer r.Close()

	// Create destination directory
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Extract files
	for _, f := range r.File {
		// Construct full path
		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(path, f.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", path, err)
			}
			continue
		}

		// Create parent directories if needed
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %v", path, err)
		}

		// Open file in zip
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s in zip: %v", f.Name, err)
		}

		// Create destination file
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create file %s: %v", path, err)
		}

		// Copy file contents
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return fmt.Errorf("failed to copy file %s: %v", path, err)
		}
	}

	return nil
}

// readResourcePackManifest reads resource pack manifest.json
func readResourcePackManifest(path string) (models.ResourcePackManifest, error) {
	var manifest models.ResourcePackManifest

	data, err := os.ReadFile(path)
	if err != nil {
		return manifest, err
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return manifest, err
	}

	return manifest, nil
}

// readWorldResourcePacks reads world_resource_packs.json
func readWorldResourcePacks(path string) ([]models.WorldResourcePack, error) {
	var packs []models.WorldResourcePack

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return packs, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &packs); err != nil {
		return nil, err
	}

	return packs, nil
}

// writeWorldResourcePacks writes world_resource_packs.json
func writeWorldResourcePacks(path string, packs []models.WorldResourcePack) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(packs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}