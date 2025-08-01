package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"minecraft-easyserver/models"
	"minecraft-easyserver/utils"
)

// System resource packs that should not be managed
var systemResourcePacks = map[string]bool{
	"vanilla":   true,
	"editor":    true,
	"chemistry": true,
}

// isSystemResourcePack checks if a resource pack is a system pack
func isSystemResourcePack(folderName string) bool {
	return systemResourcePacks[strings.ToLower(folderName)]
}

// ResourcePackService resource pack service
type ResourcePackService struct{}

// NewResourcePackService creates a new resource pack service instance
func NewResourcePackService() *ResourcePackService {
	return &ResourcePackService{}
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
	worldPath := filepath.Join(bedrockPath, "worlds", config.LevelName)
	worldResourcePacksPath := filepath.Join(worldPath, "world_resource_packs.json")
	if data, err := os.ReadFile(worldResourcePacksPath); err == nil {
		var worldPacks []models.WorldResourcePack
		if json.Unmarshal(data, &worldPacks) == nil {
			for _, pack := range worldPacks {
				activePackUUIDs[pack.PackID] = true
			}
		}
	}

	// If resource_packs directory doesn't exist, return empty list
	if _, err := os.Stat(resourcePacksPath); os.IsNotExist(err) {
		return packs, nil
	}

	entries, err := os.ReadDir(resourcePacksPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Skip system resource packs
			if isSystemResourcePack(entry.Name()) {
				continue
			}

			// Read manifest.json
			manifestPath := filepath.Join(resourcePacksPath, entry.Name(), "manifest.json")
			if data, err := os.ReadFile(manifestPath); err == nil {
				var manifest models.ResourcePackManifest
				if json.Unmarshal(data, &manifest) == nil {
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
		}
	}

	return packs, nil
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
	if err := utils.ExtractZip(zipPath, extractPath); err != nil {
		return nil, fmt.Errorf("failed to extract resource pack: %v", err)
	}

	// Clean up temporary file
	os.Remove(zipPath)

	// Read manifest.json to get pack information
	manifestPath := filepath.Join(extractPath, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		// Clean up extracted files if manifest is missing
		os.RemoveAll(extractPath)
		return nil, fmt.Errorf("manifest.json not found in resource pack")
	}

	var manifest models.ResourcePackManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		// Clean up extracted files if manifest is invalid
		os.RemoveAll(extractPath)
		return nil, fmt.Errorf("invalid manifest.json: %v", err)
	}

	packInfo := &models.ResourcePackInfo{
		Name:        manifest.Header.Name,
		UUID:        manifest.Header.UUID,
		Version:     manifest.Header.Version,
		Description: manifest.Header.Description,
		FolderName:  strings.TrimSuffix(fileName, filepath.Ext(fileName)),
		Active:      false,
	}

	return packInfo, nil
}

// ActivateResourcePack activates resource pack
func (r *ResourcePackService) ActivateResourcePack(packUUID string) error {
	// Find resource pack by UUID
	packs, err := r.GetResourcePacks()
	if err != nil {
		return err
	}

	var targetPack *models.ResourcePackInfo
	for _, pack := range packs {
		if pack.UUID == packUUID {
			targetPack = &pack
			break
		}
	}

	if targetPack == nil {
		return fmt.Errorf("resource pack not found: %s", packUUID)
	}

	// Check if it's a system resource pack
	if isSystemResourcePack(targetPack.FolderName) {
		return fmt.Errorf("cannot activate system resource pack: %s", targetPack.Name)
	}

	if targetPack.Active {
		return fmt.Errorf("resource pack already activated: %s", targetPack.Name)
	}

	// Get current world
	configPath := filepath.Join(bedrockPath, "server.properties")
	config, err := readServerProperties(configPath)
	if err != nil {
		return err
	}

	// Read current world resource packs
	worldPath := filepath.Join(bedrockPath, "worlds", config.LevelName)
	worldResourcePacksPath := filepath.Join(worldPath, "world_resource_packs.json")

	var worldPacks []models.WorldResourcePack
	if data, err := os.ReadFile(worldResourcePacksPath); err == nil {
		json.Unmarshal(data, &worldPacks)
	}

	// Add new resource pack
	newPack := models.WorldResourcePack{
		PackID:  packUUID,
		Version: targetPack.Version,
	}
	worldPacks = append(worldPacks, newPack)

	// Save updated resource packs
	data, err := json.MarshalIndent(worldPacks, "", "  ")
	if err != nil {
		return err
	}

	// Ensure world directory exists
	os.MkdirAll(worldPath, 0755)

	return os.WriteFile(worldResourcePacksPath, data, 0644)
}

// DeactivateResourcePack deactivates resource pack
func (r *ResourcePackService) DeactivateResourcePack(packUUID string) error {
	// Find resource pack by UUID to check if it's a system pack
	packs, err := r.GetResourcePacks()
	if err != nil {
		return err
	}

	var targetPack *models.ResourcePackInfo
	for _, pack := range packs {
		if pack.UUID == packUUID {
			targetPack = &pack
			break
		}
	}

	// Check if it's a system resource pack
	if targetPack != nil && isSystemResourcePack(targetPack.FolderName) {
		return fmt.Errorf("cannot deactivate system resource pack: %s", targetPack.Name)
	}

	// Get current world
	configPath := filepath.Join(bedrockPath, "server.properties")
	config, err := readServerProperties(configPath)
	if err != nil {
		return err
	}

	// Read current world resource packs
	worldPath := filepath.Join(bedrockPath, "worlds", config.LevelName)
	worldResourcePacksPath := filepath.Join(worldPath, "world_resource_packs.json")

	var worldPacks []models.WorldResourcePack
	if data, err := os.ReadFile(worldResourcePacksPath); err == nil {
		json.Unmarshal(data, &worldPacks)
	}

	// Find and remove resource pack
	found := false
	for i, pack := range worldPacks {
		if pack.PackID == packUUID {
			worldPacks = append(worldPacks[:i], worldPacks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("resource pack not activated: %s", packUUID)
	}

	// Save updated resource packs
	data, err := json.MarshalIndent(worldPacks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(worldResourcePacksPath, data, 0644)
}

// DeleteResourcePack deletes resource pack
func (r *ResourcePackService) DeleteResourcePack(packUUID string) error {
	// Find resource pack by UUID
	packs, err := r.GetResourcePacks()
	if err != nil {
		return err
	}

	var targetPack *models.ResourcePackInfo
	for _, pack := range packs {
		if pack.UUID == packUUID {
			targetPack = &pack
			break
		}
	}

	if targetPack == nil {
		return fmt.Errorf("resource pack not found: %s", packUUID)
	}

	// Check if it's a system resource pack
	if isSystemResourcePack(targetPack.FolderName) {
		return fmt.Errorf("cannot delete system resource pack: %s", targetPack.Name)
	}

	// First deactivate if active
	r.DeactivateResourcePack(packUUID) // Ignore error as pack might not be active

	// Delete resource pack directory
	resourcePacksPath := filepath.Join(bedrockPath, "resource_packs")
	packPath := filepath.Join(resourcePacksPath, targetPack.FolderName)

	return os.RemoveAll(packPath)
}