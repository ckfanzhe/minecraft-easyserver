package services

import (
	"os"
	"path/filepath"
	"testing"

	"bedrock-easyserver/models"
)

func TestConfigService(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()
	bedrockPath = tempDir

	// Create test configuration file
	configPath := filepath.Join(tempDir, "server.properties")
	testConfig := `# Minecraft Bedrock Server Configuration
server-name=Test Server
gamemode=survival
difficulty=normal
max-players=10
server-port=19132
allow-cheats=false
allow-list=true
online-mode=true
level-name=Bedrock level
default-player-permission-level=member
`

	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test configuration file: %v", err)
	}

	configService := NewConfigService()

	t.Run("GetConfig", func(t *testing.T) {
		config, err := configService.GetConfig()
		if err != nil {
			t.Fatalf("Failed to get configuration: %v", err)
		}

		if config.ServerName != "Test Server" {
			t.Errorf("Expected server name 'Test Server', got '%s'", config.ServerName)
		}

		if config.Gamemode != "survival" {
			t.Errorf("Expected gamemode 'survival', got '%s'", config.Gamemode)
		}

		if config.MaxPlayers != 10 {
			t.Errorf("Expected max players 10, got %d", config.MaxPlayers)
		}

		if !config.AllowList {
			t.Error("Expected allowlist to be enabled")
		}
	})

	t.Run("UpdateConfig", func(t *testing.T) {
		newConfig := models.ServerConfig{
			ServerName:              "Updated Server",
			Gamemode:                "creative",
			Difficulty:              "easy",
			MaxPlayers:              20,
			ServerPort:              19133,
			AllowCheats:             true,
			AllowList:               false,
			OnlineMode:              false,
			LevelName:               "New World",
			DefaultPlayerPermission: "operator",
		}

		err := configService.UpdateConfig(newConfig)
		if err != nil {
			t.Fatalf("Failed to update configuration: %v", err)
		}

		// Verify configuration is saved correctly
		savedConfig, err := configService.GetConfig()
		if err != nil {
			t.Fatalf("Failed to read saved configuration: %v", err)
		}

		if savedConfig.ServerName != "Updated Server" {
			t.Errorf("Expected server name 'Updated Server', got '%s'", savedConfig.ServerName)
		}

		if savedConfig.Gamemode != "creative" {
			t.Errorf("Expected gamemode 'creative', got '%s'", savedConfig.Gamemode)
		}

		if savedConfig.MaxPlayers != 20 {
			t.Errorf("Expected max players 20, got %d", savedConfig.MaxPlayers)
		}
	})
}

func TestAllowlistService(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()
	bedrockPath = tempDir

	allowlistService := NewAllowlistService()

	t.Run("AddToAllowlist", func(t *testing.T) {
		err := allowlistService.AddToAllowlist("TestPlayer1")
		if err != nil {
			t.Fatalf("Failed to add player to allowlist: %v", err)
		}

		err = allowlistService.AddToAllowlist("TestPlayer2")
		if err != nil {
			t.Fatalf("Failed to add player to allowlist: %v", err)
		}

		// Test duplicate addition
		err = allowlistService.AddToAllowlist("TestPlayer1")
		if err == nil {
			t.Error("Expected error when adding duplicate player")
		}
	})

	t.Run("GetAllowlist", func(t *testing.T) {
		allowlist, err := allowlistService.GetAllowlist()
		if err != nil {
			t.Fatalf("Failed to get allowlist: %v", err)
		}

		if len(allowlist) != 2 {
			t.Errorf("Expected 2 players in allowlist, got %d", len(allowlist))
		}

		expectedPlayers := map[string]bool{
			"TestPlayer1": true,
			"TestPlayer2": true,
		}

		for _, player := range allowlist {
			if !expectedPlayers[player] {
				t.Errorf("Unexpected player in allowlist: %s", player)
			}
		}
	})

	t.Run("RemoveFromAllowlist", func(t *testing.T) {
		err := allowlistService.RemoveFromAllowlist("TestPlayer1")
		if err != nil {
			t.Fatalf("Failed to remove player from allowlist: %v", err)
		}

		allowlist, err := allowlistService.GetAllowlist()
		if err != nil {
			t.Fatalf("Failed to get allowlist: %v", err)
		}

		if len(allowlist) != 1 {
			t.Errorf("Expected 1 player in allowlist, got %d", len(allowlist))
		}

		if allowlist[0] != "TestPlayer2" {
			t.Errorf("Expected remaining player 'TestPlayer2', got '%s'", allowlist[0])
		}

		// Test removing non-existent player
		err = allowlistService.RemoveFromAllowlist("NonExistentPlayer")
		if err == nil {
			t.Error("Expected error when removing non-existent player")
		}
	})

	t.Run("EmptyName", func(t *testing.T) {
		err := allowlistService.AddToAllowlist("")
		if err == nil {
			t.Error("Expected error when adding empty name")
		}

		err = allowlistService.RemoveFromAllowlist("")
		if err == nil {
			t.Error("Expected error when removing empty name")
		}
	})
}

func TestPermissionService(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()
	bedrockPath = tempDir

	permissionService := NewPermissionService()

	t.Run("UpdatePermission", func(t *testing.T) {
		err := permissionService.UpdatePermission("TestPlayer1", "operator")
		if err != nil {
			t.Fatalf("Failed to update permission: %v", err)
		}

		err = permissionService.UpdatePermission("TestPlayer2", "member")
		if err != nil {
			t.Fatalf("Failed to update permission: %v", err)
		}

		// Test invalid permission level
		err = permissionService.UpdatePermission("TestPlayer3", "invalid")
		if err == nil {
			t.Error("Expected error when setting invalid permission level")
		}
	})

	t.Run("GetPermissions", func(t *testing.T) {
		permissions, err := permissionService.GetPermissions()
		if err != nil {
			t.Fatalf("Failed to get permissions: %v", err)
		}

		if len(permissions) != 2 {
			t.Errorf("Expected 2 permission entries, got %d", len(permissions))
		}

		// Verify permission content
		foundPlayer1 := false
		foundPlayer2 := false

		for _, perm := range permissions {
			name, ok := perm["name"].(string)
			if !ok {
				t.Error("Permission entry missing name field")
				continue
			}

			level, ok := perm["level"].(string)
			if !ok {
				t.Error("Permission entry missing level field")
				continue
			}

			if name == "TestPlayer1" && level == "operator" {
				foundPlayer1 = true
			} else if name == "TestPlayer2" && level == "member" {
				foundPlayer2 = true
			}
		}

		if !foundPlayer1 {
			t.Error("TestPlayer1 operator permission not found")
		}

		if !foundPlayer2 {
			t.Error("TestPlayer2 member permission not found")
		}
	})

	t.Run("RemovePermission", func(t *testing.T) {
		err := permissionService.RemovePermission("TestPlayer1")
		if err != nil {
			t.Fatalf("Failed to remove permission: %v", err)
		}

		permissions, err := permissionService.GetPermissions()
		if err != nil {
			t.Fatalf("Failed to get permissions: %v", err)
		}

		if len(permissions) != 1 {
			t.Errorf("Expected 1 permission entry, got %d", len(permissions))
		}

		// Test removing non-existent permission
		err = permissionService.RemovePermission("NonExistentPlayer")
		if err == nil {
			t.Error("Expected error when removing non-existent permission")
		}
	})

	t.Run("EmptyName", func(t *testing.T) {
		err := permissionService.UpdatePermission("", "operator")
		if err == nil {
			t.Error("Expected error when setting permission for empty name")
		}

		err = permissionService.RemovePermission("")
		if err == nil {
			t.Error("Expected error when removing permission for empty name")
		}
	})
}

func TestWorldService(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()
	bedrockPath = tempDir

	// Create test world directory
	worldsPath := filepath.Join(tempDir, "worlds")
	os.MkdirAll(worldsPath, 0755)

	// Create test worlds
	world1Path := filepath.Join(worldsPath, "TestWorld1")
	world2Path := filepath.Join(worldsPath, "TestWorld2")
	os.MkdirAll(world1Path, 0755)
	os.MkdirAll(world2Path, 0755)

	// Create test configuration file
	configPath := filepath.Join(tempDir, "server.properties")
	testConfig := `level-name=TestWorld1`
	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test configuration file: %v", err)
	}

	worldService := NewWorldService()

	t.Run("GetWorlds", func(t *testing.T) {
		worlds, err := worldService.GetWorlds()
		if err != nil {
			t.Fatalf("Failed to get world list: %v", err)
		}

		if len(worlds) != 2 {
			t.Errorf("Expected 2 worlds, got %d", len(worlds))
		}

		// Verify active status
		foundActiveWorld := false
		for _, world := range worlds {
			if world.Name == "TestWorld1" && world.Active {
				foundActiveWorld = true
			}
		}

		if !foundActiveWorld {
			t.Error("Active world TestWorld1 not found")
		}
	})

	t.Run("ActivateWorld", func(t *testing.T) {
		err := worldService.ActivateWorld("TestWorld2")
		if err != nil {
			t.Fatalf("Failed to activate world: %v", err)
		}

		worlds, err := worldService.GetWorlds()
		if err != nil {
			t.Fatalf("Failed to get world list: %v", err)
		}

		// Verify new active status
		foundActiveWorld := false
		for _, world := range worlds {
			if world.Name == "TestWorld2" && world.Active {
				foundActiveWorld = true
			}
		}

		if !foundActiveWorld {
			t.Error("TestWorld2 was not activated correctly")
		}
	})

	t.Run("DeleteWorld", func(t *testing.T) {
		err := worldService.DeleteWorld("TestWorld1")
		if err != nil {
			t.Fatalf("Failed to delete world: %v", err)
		}

		worlds, err := worldService.GetWorlds()
		if err != nil {
			t.Fatalf("Failed to get world list: %v", err)
		}

		if len(worlds) != 1 {
			t.Errorf("Expected 1 world, got %d", len(worlds))
		}

		if worlds[0].Name != "TestWorld2" {
			t.Errorf("Expected remaining world 'TestWorld2', got '%s'", worlds[0].Name)
		}
	})

	t.Run("EmptyWorldName", func(t *testing.T) {
		err := worldService.DeleteWorld("")
		if err == nil {
			t.Error("Expected error when deleting world with empty name")
		}

		err = worldService.ActivateWorld("")
		if err == nil {
			t.Error("Expected error when activating world with empty name")
		}
	})
}