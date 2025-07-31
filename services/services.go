package services

import (
	"bufio"
	"encoding/json"
	"fmt"
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