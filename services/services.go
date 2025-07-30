package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"bedrock-easyserver/models"
)

var (
	serverProcess *exec.Cmd
	serverMutex   sync.Mutex
	bedrockPath   string
)

// InitBedrockPath 初始化bedrock路径
func InitBedrockPath(path string) error {
	if path == "" {
		return fmt.Errorf("bedrock路径不能为空")
	}
	
	// 如果是相对路径，转换为绝对路径
	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		path = filepath.Join(wd, path)
	}
	
	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("bedrock路径不存在: %s", path)
	}
	
	bedrockPath = path
	return nil
}

// SetBedrockPath 设置bedrock路径（主要用于测试）
func SetBedrockPath(path string) {
	bedrockPath = path
}

// ServerService 服务器服务
type ServerService struct{}

// NewServerService 创建新的服务器服务实例
func NewServerService() *ServerService {
	return &ServerService{}
}

// GetStatus 获取服务器状态
func (s *ServerService) GetStatus() models.ServerStatus {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverProcess == nil || serverProcess.Process == nil {
		return models.ServerStatus{
			Status:  "stopped",
			Message: "服务器未运行",
		}
	}

	// 检查进程是否仍在运行
	process, err := os.FindProcess(serverProcess.Process.Pid)
	if err != nil {
		serverProcess = nil
		return models.ServerStatus{
			Status:  "stopped",
			Message: "服务器未运行",
		}
	}

	// 在Windows上，简单检查进程是否存在
	// 如果进程已经结束，FindProcess仍然会返回一个Process对象
	// 我们可以尝试发送信号0来检查进程是否真的在运行
	if process != nil {
		return models.ServerStatus{
			Status:  "running",
			Message: "服务器正在运行",
			PID:     serverProcess.Process.Pid,
		}
	}

	serverProcess = nil
	return models.ServerStatus{
		Status:  "stopped",
		Message: "服务器未运行",
	}
}

// Start 启动服务器
func (s *ServerService) Start() error {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverProcess != nil && serverProcess.Process != nil {
		return fmt.Errorf("服务器已在运行")
	}

	exePath := filepath.Join(bedrockPath, "bedrock_server.exe")
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		return fmt.Errorf("找不到bedrock_server.exe文件")
	}

	serverProcess = exec.Command(exePath)
	serverProcess.Dir = bedrockPath

	if err := serverProcess.Start(); err != nil {
		return fmt.Errorf("启动服务器失败: %v", err)
	}

	return nil
}

// Stop 停止服务器
func (s *ServerService) Stop() error {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverProcess == nil || serverProcess.Process == nil {
		return fmt.Errorf("服务器未运行")
	}

	if err := serverProcess.Process.Kill(); err != nil {
		return fmt.Errorf("停止服务器失败: %v", err)
	}

	serverProcess.Wait()
	serverProcess = nil
	return nil
}

// Restart 重启服务器
func (s *ServerService) Restart() error {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	// 先停止
	if serverProcess != nil && serverProcess.Process != nil {
		serverProcess.Process.Kill()
		serverProcess.Wait()
		serverProcess = nil
	}

	// 等待一秒
	time.Sleep(time.Second)

	// 重新启动
	exePath := filepath.Join(bedrockPath, "bedrock_server.exe")
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		return fmt.Errorf("找不到bedrock_server.exe文件")
	}

	serverProcess = exec.Command(exePath)
	serverProcess.Dir = bedrockPath

	if err := serverProcess.Start(); err != nil {
		return fmt.Errorf("重启服务器失败: %v", err)
	}

	return nil
}

// ConfigService 配置服务
type ConfigService struct{}

// NewConfigService 创建新的配置服务实例
func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// GetConfig 获取服务器配置
func (c *ConfigService) GetConfig() (models.ServerConfig, error) {
	configPath := filepath.Join(bedrockPath, "server.properties")
	return readServerProperties(configPath)
}

// UpdateConfig 更新服务器配置
func (c *ConfigService) UpdateConfig(config models.ServerConfig) error {
	configPath := filepath.Join(bedrockPath, "server.properties")
	return writeServerProperties(configPath, config)
}

// 读取server.properties
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
			if i, err := strconv.Atoi(value); err == nil {
				config.MaxPlayers = i
			}
		case "server-port":
			if i, err := strconv.Atoi(value); err == nil {
				config.ServerPort = i
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

// 写入server.properties
func writeServerProperties(path string, config models.ServerConfig) error {
	// 读取原文件保持注释
	originalLines := []string{}
	if file, err := os.Open(path); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			originalLines = append(originalLines, scanner.Text())
		}
		file.Close()
	}

	// 创建新文件
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// 写入配置
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

	// 处理每一行
	for _, line := range originalLines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
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
			writer.WriteString(key + "=" + newValue + "\n")
			delete(configMap, key)
		} else {
			writer.WriteString(line + "\n")
		}
	}

	// 添加任何新的配置项
	for key, value := range configMap {
		writer.WriteString(key + "=" + value + "\n")
	}

	return nil
}

// AllowlistService 白名单服务
type AllowlistService struct{}

// NewAllowlistService 创建新的白名单服务实例
func NewAllowlistService() *AllowlistService {
	return &AllowlistService{}
}

// GetAllowlist 获取白名单
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

// AddToAllowlist 添加到白名单
func (a *AllowlistService) AddToAllowlist(name string) error {
	if name == "" {
		return fmt.Errorf("玩家名称不能为空")
	}

	allowlistPath := filepath.Join(bedrockPath, "allowlist.json")
	allowlist, err := readAllowlist(allowlistPath)
	if err != nil {
		allowlist = []models.AllowlistEntry{}
	}

	// 检查是否已存在
	for _, entry := range allowlist {
		if entry.Name == name {
			return fmt.Errorf("玩家已在白名单中")
		}
	}

	// 添加新条目
	newEntry := models.AllowlistEntry{
		Name:               name,
		IgnoresPlayerLimit: false,
	}
	allowlist = append(allowlist, newEntry)

	return writeAllowlist(allowlistPath, allowlist)
}

// RemoveFromAllowlist 从白名单移除
func (a *AllowlistService) RemoveFromAllowlist(name string) error {
	if name == "" {
		return fmt.Errorf("玩家名称不能为空")
	}

	allowlistPath := filepath.Join(bedrockPath, "allowlist.json")
	allowlist, err := readAllowlist(allowlistPath)
	if err != nil {
		return err
	}

	// 移除条目
	var newAllowlist []models.AllowlistEntry
	found := false
	for _, entry := range allowlist {
		if entry.Name != name {
			newAllowlist = append(newAllowlist, entry)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("玩家不在白名单中")
	}

	return writeAllowlist(allowlistPath, newAllowlist)
}

// 读取allowlist.json
func readAllowlist(path string) ([]models.AllowlistEntry, error) {
	var allowlist []models.AllowlistEntry

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return allowlist, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &allowlist); err != nil {
		return nil, err
	}

	return allowlist, nil
}

// 写入allowlist.json
func writeAllowlist(path string, allowlist []models.AllowlistEntry) error {
	data, err := json.MarshalIndent(allowlist, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// PermissionService 权限服务
type PermissionService struct{}

// NewPermissionService 创建新的权限服务实例
func NewPermissionService() *PermissionService {
	return &PermissionService{}
}

// GetPermissions 获取权限
func (p *PermissionService) GetPermissions() ([]map[string]interface{}, error) {
	permissionsPath := filepath.Join(bedrockPath, "permissions.json")
	return readPermissions(permissionsPath)
}

// UpdatePermission 更新权限
func (p *PermissionService) UpdatePermission(name, level string) error {
	if name == "" {
		return fmt.Errorf("玩家名称不能为空")
	}

	validLevels := map[string]bool{
		"visitor":  true,
		"member":   true,
		"operator": true,
	}

	if !validLevels[level] {
		return fmt.Errorf("无效的权限级别")
	}

	permissionsPath := filepath.Join(bedrockPath, "permissions.json")
	permissions, err := readPermissions(permissionsPath)
	if err != nil {
		permissions = []map[string]interface{}{}
	}

	// 查找并更新或添加权限
	found := false
	for i, perm := range permissions {
		if playerName, ok := perm["name"].(string); ok && playerName == name {
			permissions[i]["level"] = level
			found = true
			break
		}
	}

	if !found {
		newPerm := map[string]interface{}{
			"name":  name,
			"level": level,
		}
		permissions = append(permissions, newPerm)
	}

	return writePermissions(permissionsPath, permissions)
}

// RemovePermission 移除权限
func (p *PermissionService) RemovePermission(name string) error {
	if name == "" {
		return fmt.Errorf("玩家名称不能为空")
	}

	permissionsPath := filepath.Join(bedrockPath, "permissions.json")
	permissions, err := readPermissions(permissionsPath)
	if err != nil {
		return err
	}

	// 移除权限
	var newPermissions []map[string]interface{}
	found := false
	for _, perm := range permissions {
		if playerName, ok := perm["name"].(string); ok && playerName != name {
			newPermissions = append(newPermissions, perm)
		} else if ok {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("玩家权限不存在")
	}

	return writePermissions(permissionsPath, newPermissions)
}

// 读取permissions.json
func readPermissions(path string) ([]map[string]interface{}, error) {
	var permissions []map[string]interface{}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return permissions, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

// 写入permissions.json
func writePermissions(path string, permissions []map[string]interface{}) error {
	data, err := json.MarshalIndent(permissions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// WorldService 世界服务
type WorldService struct{}

// NewWorldService 创建新的世界服务实例
func NewWorldService() *WorldService {
	return &WorldService{}
}

// GetWorlds 获取世界列表
func (w *WorldService) GetWorlds() ([]models.WorldInfo, error) {
	worldsPath := filepath.Join(bedrockPath, "worlds")
	return getWorldsList(worldsPath)
}

// DeleteWorld 删除世界
func (w *WorldService) DeleteWorld(worldName string) error {
	if worldName == "" {
		return fmt.Errorf("世界名称不能为空")
	}

	// 检查世界是否存在
	worldPath := filepath.Join(bedrockPath, "worlds", worldName)
	if _, err := os.Stat(worldPath); os.IsNotExist(err) {
		return fmt.Errorf("世界不存在: %s", worldName)
	}

	// 检查是否是当前激活的世界
	configPath := filepath.Join(bedrockPath, "server.properties")
	config, err := readServerProperties(configPath)
	isActiveWorld := false
	if err == nil && config.LevelName == worldName {
		isActiveWorld = true
	}

	// 删除世界文件夹
	if err := os.RemoveAll(worldPath); err != nil {
		return fmt.Errorf("删除世界文件失败: %v", err)
	}

	// 如果删除的是当前激活的世界，需要更新配置文件
	if isActiveWorld {
		// 获取剩余的世界列表
		worldsPath := filepath.Join(bedrockPath, "worlds")
		remainingWorlds, err := getWorldsList(worldsPath)
		if err != nil {
			return fmt.Errorf("获取剩余世界列表失败: %v", err)
		}

		// 如果还有其他世界，激活第一个；否则设置为默认世界名
		if len(remainingWorlds) > 0 {
			config.LevelName = remainingWorlds[0].Name
		} else {
			// 没有其他世界时，设置为默认世界名
			config.LevelName = "Bedrock level"
		}

		// 更新配置文件
		if err := writeServerProperties(configPath, config); err != nil {
			return fmt.Errorf("更新配置文件失败: %v", err)
		}
	}

	return nil
}

// ActivateWorld 激活世界
func (w *WorldService) ActivateWorld(worldName string) error {
	if worldName == "" {
		return fmt.Errorf("世界名称不能为空")
	}

	// 更新server.properties中的level-name
	configPath := filepath.Join(bedrockPath, "server.properties")
	config, err := readServerProperties(configPath)
	if err != nil {
		return err
	}

	config.LevelName = worldName
	return writeServerProperties(configPath, config)
}

// 获取世界列表
func getWorldsList(worldsPath string) ([]models.WorldInfo, error) {
	var worlds []models.WorldInfo

	// 读取当前激活的世界
	configPath := filepath.Join(bedrockPath, "server.properties")
	config, err := readServerProperties(configPath)
	activeWorld := ""
	if err == nil {
		activeWorld = config.LevelName
	}

	entries, err := os.ReadDir(worldsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return worlds, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			world := models.WorldInfo{
				Name:   entry.Name(),
				Active: entry.Name() == activeWorld,
			}
			worlds = append(worlds, world)
		}
	}

	return worlds, nil
}