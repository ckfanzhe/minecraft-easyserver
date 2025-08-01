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
	allowlistPath := filepath.Join(bedrockPath, "allowlist.json")
	allowlist, err := readAllowlist(allowlistPath)
	if err != nil {
		return err
	}

	// Check if already exists
	for _, entry := range allowlist {
		if entry.Name == name {
			return fmt.Errorf("player %s is already in allowlist", name)
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