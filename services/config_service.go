package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"minecraft-easyserver/models"
)

// ConfigService configuration service
type ConfigService struct{}

// NewConfigService creates a new configuration service instance
func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// GetConfig gets server configuration
func (c *ConfigService) GetConfig() (models.ServerConfig, error) {
	// If no server version is active, return default configuration
	if bedrockPath == "" {
		return getDefaultServerConfig(), nil
	}
	
	configPath := filepath.Join(bedrockPath, "server.properties")
	
	// If server.properties doesn't exist, return default configuration
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return getDefaultServerConfig(), nil
	}
	
	return readServerProperties(configPath)
}

// getDefaultServerConfig returns default server configuration
func getDefaultServerConfig() models.ServerConfig {
	return models.ServerConfig{
		ServerName:               "Bedrock Server",
		Gamemode:                 "survival",
		Difficulty:               "easy",
		MaxPlayers:               10,
		ServerPort:               19132,
		AllowCheats:              false,
		AllowList:                false,
		OnlineMode:               true,
		LevelName:                "Bedrock level",
		DefaultPlayerPermission:  "member",
	}
}

// UpdateConfig updates server configuration
func (c *ConfigService) UpdateConfig(config models.ServerConfig) error {
	// If no server version is active, return error
	if bedrockPath == "" {
		return fmt.Errorf("no server version is currently active. Please download and activate a server version first")
	}
	
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
			if val, err := strconv.Atoi(value); err == nil {
				config.MaxPlayers = val
			}
		case "server-port":
			if val, err := strconv.Atoi(value); err == nil {
				config.ServerPort = val
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
	// Create configuration map
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

	// Read original file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()

	// Write new file
	writer, err := os.Create(path)
	if err != nil {
		return err
	}
	defer writer.Close()

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
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
			writer.WriteString(fmt.Sprintf("%s=%s\n", key, newValue))
		} else {
			writer.WriteString(line + "\n")
		}
	}

	return nil
}