package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"bedrock-easyserver/models"
)

func TestResourcePackService(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "bedrock_test")
	if err != nil {
		t.Fatal("Failed to create temporary directory:", err)
	}
	defer os.RemoveAll(tempDir)

	// Set bedrock path
	SetBedrockPath(tempDir)

	// Create test server.properties first
	serverPropsPath := filepath.Join(tempDir, "server.properties")
	serverPropsContent := "level-name=test_world\n"
	if err := os.WriteFile(serverPropsPath, []byte(serverPropsContent), 0644); err != nil {
		t.Fatal("Failed to write server.properties:", err)
	}

	// Create test world directory
	worldsPath := filepath.Join(tempDir, "worlds", "test_world")
	if err := os.MkdirAll(worldsPath, 0755); err != nil {
		t.Fatal("Failed to create world directory:", err)
	}

	// Create resource pack service
	service := NewResourcePackService()

	// Test GetResourcePacks with empty directory
	t.Run("GetResourcePacks_EmptyDirectory", func(t *testing.T) {
		packs, err := service.GetResourcePacks()
		if err != nil {
			t.Fatal("Failed to get resource packs:", err)
		}
		if len(packs) != 0 {
			t.Errorf("Expected 0 resource packs, got %d", len(packs))
		}
	})

	// Create test resource pack directory structure
	resourcePacksPath := filepath.Join(tempDir, "resource_packs")
	testPackPath := filepath.Join(resourcePacksPath, "test_pack")
	if err := os.MkdirAll(testPackPath, 0755); err != nil {
		t.Fatal("Failed to create test pack directory:", err)
	}

	// Create test manifest.json
	manifest := models.ResourcePackManifest{
		FormatVersion: 2,
		Header: models.ResourcePackHeader{
			Description: "Test Resource Pack",
			Name:        "Test Pack",
			UUID:        "12345678-1234-1234-1234-123456789012",
			Version:     [3]int{1, 0, 0},
		},
		Modules: []models.ResourcePackModule{
			{
				Description: "Test Module",
				Type:        "resources",
				UUID:        "87654321-4321-4321-4321-210987654321",
				Version:     [3]int{1, 0, 0},
			},
		},
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal("Failed to marshal manifest:", err)
	}

	manifestPath := filepath.Join(testPackPath, "manifest.json")
	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		t.Fatal("Failed to write manifest file:", err)
	}

	// Test GetResourcePacks with test pack
	t.Run("GetResourcePacks_WithTestPack", func(t *testing.T) {
		packs, err := service.GetResourcePacks()
		if err != nil {
			t.Fatal("Failed to get resource packs:", err)
		}
		if len(packs) != 1 {
			t.Errorf("Expected 1 resource pack, got %d", len(packs))
		}
		if packs[0].Name != "Test Pack" {
			t.Errorf("Expected pack name 'Test Pack', got '%s'", packs[0].Name)
		}
		if packs[0].UUID != "12345678-1234-1234-1234-123456789012" {
			t.Errorf("Expected pack UUID '12345678-1234-1234-1234-123456789012', got '%s'", packs[0].UUID)
		}
		if packs[0].Active {
			t.Error("Expected pack to be inactive initially")
		}
	})

	// Test ActivateResourcePack
	t.Run("ActivateResourcePack", func(t *testing.T) {
		err := service.ActivateResourcePack("12345678-1234-1234-1234-123456789012")
		if err != nil {
			t.Fatal("Failed to activate resource pack:", err)
		}

		// Check if world_resource_packs.json was created
		worldResourcePacksPath := filepath.Join(worldsPath, "world_resource_packs.json")
		if _, err := os.Stat(worldResourcePacksPath); os.IsNotExist(err) {
			t.Error("world_resource_packs.json was not created")
		}

		// Read and verify content
		data, err := os.ReadFile(worldResourcePacksPath)
		if err != nil {
			t.Fatal("Failed to read world_resource_packs.json:", err)
		}

		var activePacks []models.WorldResourcePack
		if err := json.Unmarshal(data, &activePacks); err != nil {
			t.Fatal("Failed to unmarshal world_resource_packs.json:", err)
		}

		if len(activePacks) != 1 {
			t.Errorf("Expected 1 active pack, got %d", len(activePacks))
		}
		if activePacks[0].PackID != "12345678-1234-1234-1234-123456789012" {
			t.Errorf("Expected pack ID '12345678-1234-1234-1234-123456789012', got '%s'", activePacks[0].PackID)
		}
	})

	// Test GetResourcePacks after activation
	t.Run("GetResourcePacks_AfterActivation", func(t *testing.T) {
		packs, err := service.GetResourcePacks()
		if err != nil {
			t.Fatal("Failed to get resource packs:", err)
		}
		if len(packs) != 1 {
			t.Errorf("Expected 1 resource pack, got %d", len(packs))
		}
		if !packs[0].Active {
			t.Error("Expected pack to be active after activation")
		}
	})

	// Test ActivateResourcePack with already activated pack
	t.Run("ActivateResourcePack_AlreadyActivated", func(t *testing.T) {
		err := service.ActivateResourcePack("12345678-1234-1234-1234-123456789012")
		if err == nil {
			t.Error("Expected error when activating already activated pack")
		}
	})

	// Test DeactivateResourcePack
	t.Run("DeactivateResourcePack", func(t *testing.T) {
		err := service.DeactivateResourcePack("12345678-1234-1234-1234-123456789012")
		if err != nil {
			t.Fatal("Failed to deactivate resource pack:", err)
		}

		// Check if pack is removed from world_resource_packs.json
		worldResourcePacksPath := filepath.Join(worldsPath, "world_resource_packs.json")
		data, err := os.ReadFile(worldResourcePacksPath)
		if err != nil {
			t.Fatal("Failed to read world_resource_packs.json:", err)
		}

		var activePacks []models.WorldResourcePack
		if err := json.Unmarshal(data, &activePacks); err != nil {
			t.Fatal("Failed to unmarshal world_resource_packs.json:", err)
		}

		if len(activePacks) != 0 {
			t.Errorf("Expected 0 active packs, got %d", len(activePacks))
		}
	})

	// Test DeactivateResourcePack with not activated pack
	t.Run("DeactivateResourcePack_NotActivated", func(t *testing.T) {
		err := service.DeactivateResourcePack("12345678-1234-1234-1234-123456789012")
		if err == nil {
			t.Error("Expected error when deactivating not activated pack")
		}
	})

	// Test DeleteResourcePack
	t.Run("DeleteResourcePack", func(t *testing.T) {
		err := service.DeleteResourcePack("12345678-1234-1234-1234-123456789012")
		if err != nil {
			t.Fatal("Failed to delete resource pack:", err)
		}

		// Check if directory was deleted
		if _, err := os.Stat(testPackPath); !os.IsNotExist(err) {
			t.Error("Resource pack directory was not deleted")
		}
	})

	// Test GetResourcePacks after deletion
	t.Run("GetResourcePacks_AfterDeletion", func(t *testing.T) {
		packs, err := service.GetResourcePacks()
		if err != nil {
			t.Fatal("Failed to get resource packs:", err)
		}
		if len(packs) != 0 {
			t.Errorf("Expected 0 resource packs after deletion, got %d", len(packs))
		}
	})

	// Test with non-existent pack
	t.Run("ActivateResourcePack_NotFound", func(t *testing.T) {
		err := service.ActivateResourcePack("non-existent-uuid")
		if err == nil {
			t.Error("Expected error when activating non-existent pack")
		}
	})

	t.Run("DeleteResourcePack_NotFound", func(t *testing.T) {
		err := service.DeleteResourcePack("non-existent-uuid")
		if err == nil {
			t.Error("Expected error when deleting non-existent pack")
		}
	})
}

func TestResourcePackHelperFunctions(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "bedrock_test")
	if err != nil {
		t.Fatal("Failed to create temporary directory:", err)
	}
	defer os.RemoveAll(tempDir)

	// Test readResourcePackManifest
	t.Run("readResourcePackManifest", func(t *testing.T) {
		manifest := models.ResourcePackManifest{
			FormatVersion: 2,
			Header: models.ResourcePackHeader{
				Description: "Test Resource Pack",
				Name:        "Test Pack",
				UUID:        "12345678-1234-1234-1234-123456789012",
				Version:     [3]int{1, 0, 0},
			},
		}

		manifestData, err := json.MarshalIndent(manifest, "", "  ")
		if err != nil {
			t.Fatal("Failed to marshal manifest:", err)
		}

		manifestPath := filepath.Join(tempDir, "manifest.json")
		if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
			t.Fatal("Failed to write manifest file:", err)
		}

		readManifest, err := readResourcePackManifest(manifestPath)
		if err != nil {
			t.Fatal("Failed to read manifest:", err)
		}

		if readManifest.Header.Name != "Test Pack" {
			t.Errorf("Expected name 'Test Pack', got '%s'", readManifest.Header.Name)
		}
		if readManifest.Header.UUID != "12345678-1234-1234-1234-123456789012" {
			t.Errorf("Expected UUID '12345678-1234-1234-1234-123456789012', got '%s'", readManifest.Header.UUID)
		}
	})

	// Test readWorldResourcePacks and writeWorldResourcePacks
	t.Run("WorldResourcePacks_ReadWrite", func(t *testing.T) {
		packs := []models.WorldResourcePack{
			{
				PackID:  "12345678-1234-1234-1234-123456789012",
				Version: [3]int{1, 0, 0},
			},
			{
				PackID:  "87654321-4321-4321-4321-210987654321",
				Version: [3]int{2, 1, 0},
			},
		}

		worldResourcePacksPath := filepath.Join(tempDir, "world_resource_packs.json")

		// Test write
		err := writeWorldResourcePacks(worldResourcePacksPath, packs)
		if err != nil {
			t.Fatal("Failed to write world resource packs:", err)
		}

		// Test read
		readPacks, err := readWorldResourcePacks(worldResourcePacksPath)
		if err != nil {
			t.Fatal("Failed to read world resource packs:", err)
		}

		if len(readPacks) != 2 {
			t.Errorf("Expected 2 packs, got %d", len(readPacks))
		}
		if readPacks[0].PackID != "12345678-1234-1234-1234-123456789012" {
			t.Errorf("Expected first pack ID '12345678-1234-1234-1234-123456789012', got '%s'", readPacks[0].PackID)
		}
		if readPacks[1].PackID != "87654321-4321-4321-4321-210987654321" {
			t.Errorf("Expected second pack ID '87654321-4321-4321-4321-210987654321', got '%s'", readPacks[1].PackID)
		}
	})

	// Test readWorldResourcePacks with non-existent file
	t.Run("readWorldResourcePacks_NonExistent", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "non_existent.json")
		packs, err := readWorldResourcePacks(nonExistentPath)
		if err != nil {
			t.Fatal("Expected no error for non-existent file:", err)
		}
		if len(packs) != 0 {
			t.Errorf("Expected 0 packs for non-existent file, got %d", len(packs))
		}
	})
}