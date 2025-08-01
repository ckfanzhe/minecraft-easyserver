package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

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
	permissionsPath := filepath.Join(bedrockPath, "permissions.json")
	permissions, err := readPermissions(permissionsPath)
	if err != nil {
		return err
	}

	// Validate permission level
	validLevels := []string{"visitor", "member", "operator"}
	isValid := false
	for _, validLevel := range validLevels {
		if level == validLevel {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid permission level: %s. Valid levels are: visitor, member, operator", level)
	}

	// Find and update existing permission
	found := false
	for i, permission := range permissions {
		if permission["name"] == name {
			permissions[i]["permission"] = level
			found = true
			break
		}
	}

	// If not found, add new permission
	if !found {
		newPermission := map[string]interface{}{
			"name":       name,
			"permission": level,
		}
		permissions = append(permissions, newPermission)
	}

	return writePermissions(permissionsPath, permissions)
}

// RemovePermission removes permission
func (p *PermissionService) RemovePermission(name string) error {
	permissionsPath := filepath.Join(bedrockPath, "permissions.json")
	permissions, err := readPermissions(permissionsPath)
	if err != nil {
		return err
	}

	// Find and remove permission
	found := false
	for i, permission := range permissions {
		if permission["name"] == name {
			permissions = append(permissions[:i], permissions[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("permission not found for player: %s", name)
	}

	return writePermissions(permissionsPath, permissions)
}

// Read permissions.json
func readPermissions(path string) ([]map[string]interface{}, error) {
	var permissions []map[string]interface{}

	// If file doesn't exist, return empty permissions
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return permissions, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
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