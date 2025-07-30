package services

import (
	"os"
	"path/filepath"
	"testing"

	"bedrock-easyserver/models"
)

func TestConfigService(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	bedrockPath = tempDir

	// 创建测试配置文件
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
		t.Fatalf("创建测试配置文件失败: %v", err)
	}

	configService := NewConfigService()

	t.Run("GetConfig", func(t *testing.T) {
		config, err := configService.GetConfig()
		if err != nil {
			t.Fatalf("获取配置失败: %v", err)
		}

		if config.ServerName != "Test Server" {
			t.Errorf("期望服务器名称为 'Test Server', 实际为 '%s'", config.ServerName)
		}

		if config.Gamemode != "survival" {
			t.Errorf("期望游戏模式为 'survival', 实际为 '%s'", config.Gamemode)
		}

		if config.MaxPlayers != 10 {
			t.Errorf("期望最大玩家数为 10, 实际为 %d", config.MaxPlayers)
		}

		if !config.AllowList {
			t.Error("期望白名单开启")
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
			t.Fatalf("更新配置失败: %v", err)
		}

		// 验证配置是否正确保存
		savedConfig, err := configService.GetConfig()
		if err != nil {
			t.Fatalf("读取保存的配置失败: %v", err)
		}

		if savedConfig.ServerName != "Updated Server" {
			t.Errorf("期望服务器名称为 'Updated Server', 实际为 '%s'", savedConfig.ServerName)
		}

		if savedConfig.Gamemode != "creative" {
			t.Errorf("期望游戏模式为 'creative', 实际为 '%s'", savedConfig.Gamemode)
		}

		if savedConfig.MaxPlayers != 20 {
			t.Errorf("期望最大玩家数为 20, 实际为 %d", savedConfig.MaxPlayers)
		}
	})
}

func TestAllowlistService(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	bedrockPath = tempDir

	allowlistService := NewAllowlistService()

	t.Run("AddToAllowlist", func(t *testing.T) {
		err := allowlistService.AddToAllowlist("TestPlayer1")
		if err != nil {
			t.Fatalf("添加玩家到白名单失败: %v", err)
		}

		err = allowlistService.AddToAllowlist("TestPlayer2")
		if err != nil {
			t.Fatalf("添加玩家到白名单失败: %v", err)
		}

		// 测试重复添加
		err = allowlistService.AddToAllowlist("TestPlayer1")
		if err == nil {
			t.Error("期望重复添加玩家时返回错误")
		}
	})

	t.Run("GetAllowlist", func(t *testing.T) {
		allowlist, err := allowlistService.GetAllowlist()
		if err != nil {
			t.Fatalf("获取白名单失败: %v", err)
		}

		if len(allowlist) != 2 {
			t.Errorf("期望白名单有2个玩家, 实际为 %d", len(allowlist))
		}

		expectedPlayers := map[string]bool{
			"TestPlayer1": true,
			"TestPlayer2": true,
		}

		for _, player := range allowlist {
			if !expectedPlayers[player] {
				t.Errorf("意外的玩家在白名单中: %s", player)
			}
		}
	})

	t.Run("RemoveFromAllowlist", func(t *testing.T) {
		err := allowlistService.RemoveFromAllowlist("TestPlayer1")
		if err != nil {
			t.Fatalf("从白名单移除玩家失败: %v", err)
		}

		allowlist, err := allowlistService.GetAllowlist()
		if err != nil {
			t.Fatalf("获取白名单失败: %v", err)
		}

		if len(allowlist) != 1 {
			t.Errorf("期望白名单有1个玩家, 实际为 %d", len(allowlist))
		}

		if allowlist[0] != "TestPlayer2" {
			t.Errorf("期望剩余玩家为 'TestPlayer2', 实际为 '%s'", allowlist[0])
		}

		// 测试移除不存在的玩家
		err = allowlistService.RemoveFromAllowlist("NonExistentPlayer")
		if err == nil {
			t.Error("期望移除不存在的玩家时返回错误")
		}
	})

	t.Run("EmptyName", func(t *testing.T) {
		err := allowlistService.AddToAllowlist("")
		if err == nil {
			t.Error("期望添加空名称时返回错误")
		}

		err = allowlistService.RemoveFromAllowlist("")
		if err == nil {
			t.Error("期望移除空名称时返回错误")
		}
	})
}

func TestPermissionService(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	bedrockPath = tempDir

	permissionService := NewPermissionService()

	t.Run("UpdatePermission", func(t *testing.T) {
		err := permissionService.UpdatePermission("TestPlayer1", "operator")
		if err != nil {
			t.Fatalf("更新权限失败: %v", err)
		}

		err = permissionService.UpdatePermission("TestPlayer2", "member")
		if err != nil {
			t.Fatalf("更新权限失败: %v", err)
		}

		// 测试无效权限级别
		err = permissionService.UpdatePermission("TestPlayer3", "invalid")
		if err == nil {
			t.Error("期望设置无效权限级别时返回错误")
		}
	})

	t.Run("GetPermissions", func(t *testing.T) {
		permissions, err := permissionService.GetPermissions()
		if err != nil {
			t.Fatalf("获取权限失败: %v", err)
		}

		if len(permissions) != 2 {
			t.Errorf("期望有2个权限条目, 实际为 %d", len(permissions))
		}

		// 验证权限内容
		foundPlayer1 := false
		foundPlayer2 := false

		for _, perm := range permissions {
			name, ok := perm["name"].(string)
			if !ok {
				t.Error("权限条目缺少name字段")
				continue
			}

			level, ok := perm["level"].(string)
			if !ok {
				t.Error("权限条目缺少level字段")
				continue
			}

			if name == "TestPlayer1" && level == "operator" {
				foundPlayer1 = true
			} else if name == "TestPlayer2" && level == "member" {
				foundPlayer2 = true
			}
		}

		if !foundPlayer1 {
			t.Error("未找到TestPlayer1的operator权限")
		}

		if !foundPlayer2 {
			t.Error("未找到TestPlayer2的member权限")
		}
	})

	t.Run("RemovePermission", func(t *testing.T) {
		err := permissionService.RemovePermission("TestPlayer1")
		if err != nil {
			t.Fatalf("移除权限失败: %v", err)
		}

		permissions, err := permissionService.GetPermissions()
		if err != nil {
			t.Fatalf("获取权限失败: %v", err)
		}

		if len(permissions) != 1 {
			t.Errorf("期望有1个权限条目, 实际为 %d", len(permissions))
		}

		// 测试移除不存在的权限
		err = permissionService.RemovePermission("NonExistentPlayer")
		if err == nil {
			t.Error("期望移除不存在的权限时返回错误")
		}
	})

	t.Run("EmptyName", func(t *testing.T) {
		err := permissionService.UpdatePermission("", "operator")
		if err == nil {
			t.Error("期望设置空名称权限时返回错误")
		}

		err = permissionService.RemovePermission("")
		if err == nil {
			t.Error("期望移除空名称权限时返回错误")
		}
	})
}

func TestWorldService(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	bedrockPath = tempDir

	// 创建测试世界目录
	worldsPath := filepath.Join(tempDir, "worlds")
	os.MkdirAll(worldsPath, 0755)

	// 创建测试世界
	world1Path := filepath.Join(worldsPath, "TestWorld1")
	world2Path := filepath.Join(worldsPath, "TestWorld2")
	os.MkdirAll(world1Path, 0755)
	os.MkdirAll(world2Path, 0755)

	// 创建测试配置文件
	configPath := filepath.Join(tempDir, "server.properties")
	testConfig := `level-name=TestWorld1`
	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("创建测试配置文件失败: %v", err)
	}

	worldService := NewWorldService()

	t.Run("GetWorlds", func(t *testing.T) {
		worlds, err := worldService.GetWorlds()
		if err != nil {
			t.Fatalf("获取世界列表失败: %v", err)
		}

		if len(worlds) != 2 {
			t.Errorf("期望有2个世界, 实际为 %d", len(worlds))
		}

		// 验证激活状态
		foundActiveWorld := false
		for _, world := range worlds {
			if world.Name == "TestWorld1" && world.Active {
				foundActiveWorld = true
			}
		}

		if !foundActiveWorld {
			t.Error("未找到激活的世界TestWorld1")
		}
	})

	t.Run("ActivateWorld", func(t *testing.T) {
		err := worldService.ActivateWorld("TestWorld2")
		if err != nil {
			t.Fatalf("激活世界失败: %v", err)
		}

		worlds, err := worldService.GetWorlds()
		if err != nil {
			t.Fatalf("获取世界列表失败: %v", err)
		}

		// 验证新的激活状态
		foundActiveWorld := false
		for _, world := range worlds {
			if world.Name == "TestWorld2" && world.Active {
				foundActiveWorld = true
			}
		}

		if !foundActiveWorld {
			t.Error("TestWorld2未被正确激活")
		}
	})

	t.Run("DeleteWorld", func(t *testing.T) {
		err := worldService.DeleteWorld("TestWorld1")
		if err != nil {
			t.Fatalf("删除世界失败: %v", err)
		}

		worlds, err := worldService.GetWorlds()
		if err != nil {
			t.Fatalf("获取世界列表失败: %v", err)
		}

		if len(worlds) != 1 {
			t.Errorf("期望有1个世界, 实际为 %d", len(worlds))
		}

		if worlds[0].Name != "TestWorld2" {
			t.Errorf("期望剩余世界为 'TestWorld2', 实际为 '%s'", worlds[0].Name)
		}
	})

	t.Run("EmptyWorldName", func(t *testing.T) {
		err := worldService.DeleteWorld("")
		if err == nil {
			t.Error("期望删除空名称世界时返回错误")
		}

		err = worldService.ActivateWorld("")
		if err == nil {
			t.Error("期望激活空名称世界时返回错误")
		}
	})
}