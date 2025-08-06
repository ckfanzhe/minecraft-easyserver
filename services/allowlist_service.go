package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"minecraft-easyserver/models"
)

// AllowlistService allowlist service
type AllowlistService struct{}

// NewAllowlistService creates a new allowlist service instance
func NewAllowlistService() *AllowlistService {
	return &AllowlistService{}
}

// GetAllowlist gets allowlist
func (a *AllowlistService) GetAllowlist() ([]models.AllowlistEntry, error) {
	// If no server version is active, return empty list
	if bedrockPath == "" {
		return []models.AllowlistEntry{}, nil
	}
	
	allowlistPath := filepath.Join(bedrockPath, "allowlist.json")
	allowlist, err := readAllowlist(allowlistPath)
	if err != nil {
		return nil, err
	}

	return allowlist, nil
}

// AddToAllowlist adds to allowlist
func (a *AllowlistService) AddToAllowlist(entry models.AllowlistEntry) error {
	// If no server version is active, return error
	if bedrockPath == "" {
		return fmt.Errorf("no server version is currently active. Please download and activate a server version first")
	}
	
	allowlistPath := filepath.Join(bedrockPath, "allowlist.json")
	allowlist, err := readAllowlist(allowlistPath)
	if err != nil {
		return err
	}

	// Check if already exists
	for _, existingEntry := range allowlist {
		if existingEntry.Name == entry.Name {
			return fmt.Errorf("player %s is already in allowlist", entry.Name)
		}
	}

	// Add new entry
	newEntry := models.AllowlistEntry{
		Name:               entry.Name,
		IgnoresPlayerLimit: entry.IgnoresPlayerLimit,
	}
	allowlist = append(allowlist, newEntry)

	return writeAllowlist(allowlistPath, allowlist)
}

// RemoveFromAllowlist removes from allowlist
func (a *AllowlistService) RemoveFromAllowlist(name string) error {
	// If no server version is active, return error
	if bedrockPath == "" {
		return fmt.Errorf("no server version is currently active. Please download and activate a server version first")
	}
	
	allowlistPath := filepath.Join(bedrockPath, "allowlist.json")
	allowlist, err := readAllowlist(allowlistPath)
	if err != nil {
		return err
	}

	// Find and remove entry
	found := false
	for i, entry := range allowlist {
		if entry.Name == name {
			allowlist = append(allowlist[:i], allowlist[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("player %s not in allowlist", name)
	}

	return writeAllowlist(allowlistPath, allowlist)
}

// Read allowlist.json
func readAllowlist(path string) ([]models.AllowlistEntry, error) {
	var allowlist []models.AllowlistEntry

	// If file doesn't exist, return empty allowlist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return allowlist, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
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