package services

import (
	"fmt"
	"os"
	"path/filepath"

	"minecraft-easyserver/models"
)

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

	// Delete world directory
	if err := os.RemoveAll(worldPath); err != nil {
		return fmt.Errorf("failed to delete world directory: %v", err)
	}

	// If this is the active world, reset to default
	configPath := filepath.Join(bedrockPath, "server.properties")
	config, err := readServerProperties(configPath)
	if err != nil {
		return err
	}

	if config.LevelName == worldName {
		config.LevelName = "Bedrock level"
		if err := writeServerProperties(configPath, config); err != nil {
			return err
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

	// If worlds directory doesn't exist, return empty list
	if _, err := os.Stat(worldsPath); os.IsNotExist(err) {
		return worlds, nil
	}

	entries, err := os.ReadDir(worldsPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Check if it's a valid world (contains level.dat or levelname.txt)
			worldPath := filepath.Join(worldsPath, entry.Name())
			levelDatPath := filepath.Join(worldPath, "level.dat")
			levelNamePath := filepath.Join(worldPath, "levelname.txt")

			if _, err := os.Stat(levelDatPath); err == nil {
				// Valid world
				worldInfo := models.WorldInfo{
					Name:   entry.Name(),
					Active: entry.Name() == activeWorld,
				}
				worlds = append(worlds, worldInfo)
			} else if _, err := os.Stat(levelNamePath); err == nil {
				// Valid world
				worldInfo := models.WorldInfo{
					Name:   entry.Name(),
					Active: entry.Name() == activeWorld,
				}
				worlds = append(worlds, worldInfo)
			}
		}
	}

	return worlds, nil
}