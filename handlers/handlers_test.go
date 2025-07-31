package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"bedrock-easyserver/models"
	"bedrock-easyserver/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestEnvironment(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "bedrock-test-*")
	assert.NoError(t, err)

	// Set bedrock path for test environment
	services.SetBedrockPath(tempDir)

	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return tempDir
}

func TestServerHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tempDir := setupTestEnvironment(t)
	_ = tempDir

	handler := NewServerHandler()

	t.Run("GetStatus", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/status", nil)

		handler.GetStatus(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "stopped", response["status"])
	})
}

func TestConfigHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tempDir := setupTestEnvironment(t)
	
	// Create test configuration file
	configPath := filepath.Join(tempDir, "server.properties")
	testConfig := `server-name=Test Server
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
	assert.NoError(t, err)

	handler := NewConfigHandler()

	t.Run("GetConfig", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/config", nil)

		handler.GetConfig(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		config, exists := response["config"]
		assert.True(t, exists)
		
		configMap := config.(map[string]interface{})
		assert.Equal(t, "Test Server", configMap["server-name"])
		assert.Equal(t, "survival", configMap["gamemode"])
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

		requestBody := map[string]interface{}{
			"config": newConfig,
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/config", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateConfig(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Configuration saved")
	})

	t.Run("UpdateConfig_InvalidJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/config", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateConfig(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAllowlistHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tempDir := setupTestEnvironment(t)
	_ = tempDir

	handler := NewAllowlistHandler()

	t.Run("AddToAllowlist", func(t *testing.T) {
		requestBody := map[string]string{
			"name": "TestPlayer1",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/allowlist", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddToAllowlist(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Added to allowlist")
	})

	t.Run("GetAllowlist", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/allowlist", nil)

		handler.GetAllowlist(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		allowlist, exists := response["allowlist"]
		assert.True(t, exists)
		
		allowlistSlice := allowlist.([]interface{})
		assert.Len(t, allowlistSlice, 1)
		assert.Equal(t, "TestPlayer1", allowlistSlice[0])
	})

	t.Run("RemoveFromAllowlist", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("DELETE", "/api/allowlist/TestPlayer1", nil)
		c.Params = gin.Params{gin.Param{Key: "name", Value: "TestPlayer1"}}

		handler.RemoveFromAllowlist(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Removed from allowlist")
	})

	t.Run("AddToAllowlist_EmptyName", func(t *testing.T) {
		requestBody := map[string]string{
			"name": "",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/allowlist", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddToAllowlist(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("AddToAllowlist_InvalidJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/allowlist", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddToAllowlist(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPermissionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tempDir := setupTestEnvironment(t)
	_ = tempDir

	handler := NewPermissionHandler()

	t.Run("UpdatePermission", func(t *testing.T) {
		requestBody := map[string]string{
			"name":  "TestPlayer1",
			"level": "operator",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/permissions", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdatePermission(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Permission set")
	})

	t.Run("GetPermissions", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/permissions", nil)

		handler.GetPermissions(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		permissions, exists := response["permissions"]
		assert.True(t, exists)
		
		permissionsSlice := permissions.([]interface{})
		assert.Len(t, permissionsSlice, 1)
	})

	t.Run("UpdatePermission_InvalidLevel", func(t *testing.T) {
		requestBody := map[string]string{
			"name":  "TestPlayer1",
			"level": "invalid",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/permissions", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdatePermission(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("RemovePermission", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("DELETE", "/api/permissions/TestPlayer1", nil)
		c.Params = gin.Params{gin.Param{Key: "name", Value: "TestPlayer1"}}

		handler.RemovePermission(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Permission removed")
	})
}

func TestWorldHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tempDir := setupTestEnvironment(t)

	// Create test world directory
	worldsPath := filepath.Join(tempDir, "worlds")
	os.MkdirAll(worldsPath, 0755)

	// Create test world
	world1Path := filepath.Join(worldsPath, "TestWorld1")
	os.MkdirAll(world1Path, 0755)

	// Create test configuration file
	configPath := filepath.Join(tempDir, "server.properties")
	testConfig := `level-name=TestWorld1`
	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	assert.NoError(t, err)

	handler := NewWorldHandler()

	t.Run("GetWorlds", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/worlds", nil)

		handler.GetWorlds(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		worlds, exists := response["worlds"]
		assert.True(t, exists)
		
		worldsSlice := worlds.([]interface{})
		assert.Len(t, worldsSlice, 1)
	})

	t.Run("DeleteWorld", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("DELETE", "/api/worlds/TestWorld1", nil)
		c.Params = gin.Params{gin.Param{Key: "name", Value: "TestWorld1"}}

		handler.DeleteWorld(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "World deleted")
	})

	t.Run("ActivateWorld", func(t *testing.T) {
		// Recreate world for activation test
		world2Path := filepath.Join(worldsPath, "TestWorld2")
		os.MkdirAll(world2Path, 0755)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/worlds/TestWorld2/activate", nil)
		c.Params = gin.Params{gin.Param{Key: "name", Value: "TestWorld2"}}

		handler.ActivateWorld(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "World activated")
	})
}