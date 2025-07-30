package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"bedrock-easyserver/models"
	"bedrock-easyserver/services"
	"github.com/gin-gonic/gin"
)

// ServerHandler 服务器处理器
type ServerHandler struct {
	serverService *services.ServerService
}

// NewServerHandler 创建新的服务器处理器
func NewServerHandler() *ServerHandler {
	return &ServerHandler{
		serverService: services.NewServerService(),
	}
}

// GetStatus 获取服务器状态
func (h *ServerHandler) GetStatus(c *gin.Context) {
	status := h.serverService.GetStatus()
	c.JSON(200, status)
}

// StartServer 启动服务器
func (h *ServerHandler) StartServer(c *gin.Context) {
	if err := h.serverService.Start(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "服务器启动成功"})
}

// StopServer 停止服务器
func (h *ServerHandler) StopServer(c *gin.Context) {
	if err := h.serverService.Stop(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "服务器已停止"})
}

// RestartServer 重启服务器
func (h *ServerHandler) RestartServer(c *gin.Context) {
	if err := h.serverService.Restart(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "服务器重启成功"})
}

// ConfigHandler 配置处理器
type ConfigHandler struct {
	configService *services.ConfigService
}

// NewConfigHandler 创建新的配置处理器
func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		configService: services.NewConfigService(),
	}
}

// GetConfig 获取服务器配置
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	config, err := h.configService.GetConfig()
	if err != nil {
		c.JSON(500, gin.H{"error": "读取配置失败: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"config": config})
}

// UpdateConfig 更新服务器配置
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	var request struct {
		Config models.ServerConfig `json:"config"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := h.configService.UpdateConfig(request.Config); err != nil {
		c.JSON(500, gin.H{"error": "保存配置失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "配置已保存，重启服务器后生效"})
}

// AllowlistHandler 白名单处理器
type AllowlistHandler struct {
	allowlistService *services.AllowlistService
}

// NewAllowlistHandler 创建新的白名单处理器
func NewAllowlistHandler() *AllowlistHandler {
	return &AllowlistHandler{
		allowlistService: services.NewAllowlistService(),
	}
}

// GetAllowlist 获取白名单
func (h *AllowlistHandler) GetAllowlist(c *gin.Context) {
	allowlist, err := h.allowlistService.GetAllowlist()
	if err != nil {
		c.JSON(500, gin.H{"error": "读取白名单失败: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"allowlist": allowlist})
}

// AddToAllowlist 添加到白名单
func (h *AllowlistHandler) AddToAllowlist(c *gin.Context) {
	var request struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := h.allowlistService.AddToAllowlist(request.Name); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "已添加到白名单: " + request.Name})
}

// RemoveFromAllowlist 从白名单移除
func (h *AllowlistHandler) RemoveFromAllowlist(c *gin.Context) {
	name := c.Param("name")

	if err := h.allowlistService.RemoveFromAllowlist(name); err != nil {
		if strings.Contains(err.Error(), "不在白名单中") {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "保存白名单失败: " + err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "已从白名单移除: " + name})
}

// PermissionHandler 权限处理器
type PermissionHandler struct {
	permissionService *services.PermissionService
}

// NewPermissionHandler 创建新的权限处理器
func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{
		permissionService: services.NewPermissionService(),
	}
}

// GetPermissions 获取权限
func (h *PermissionHandler) GetPermissions(c *gin.Context) {
	permissions, err := h.permissionService.GetPermissions()
	if err != nil {
		c.JSON(500, gin.H{"error": "读取权限失败: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"permissions": permissions})
}

// UpdatePermission 更新权限
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	var request struct {
		Name  string `json:"name"`
		Level string `json:"level"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := h.permissionService.UpdatePermission(request.Name, request.Level); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("已设置 %s 的权限为 %s", request.Name, request.Level)})
}

// RemovePermission 移除权限
func (h *PermissionHandler) RemovePermission(c *gin.Context) {
	name := c.Param("name")

	if err := h.permissionService.RemovePermission(name); err != nil {
		if strings.Contains(err.Error(), "权限不存在") {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "保存权限失败: " + err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "已移除权限: " + name})
}

// WorldHandler 世界处理器
type WorldHandler struct {
	worldService *services.WorldService
}

// NewWorldHandler 创建新的世界处理器
func NewWorldHandler() *WorldHandler {
	return &WorldHandler{
		worldService: services.NewWorldService(),
	}
}

// GetWorlds 获取世界列表
func (h *WorldHandler) GetWorlds(c *gin.Context) {
	worlds, err := h.worldService.GetWorlds()
	if err != nil {
		c.JSON(500, gin.H{"error": "读取世界列表失败: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"worlds": worlds})
}

// extractZip 解压zip文件
func extractZip(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 确保目标目录存在
	os.MkdirAll(dest, 0755)

	// 提取文件
	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)

		// 检查路径安全性，防止目录遍历攻击
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("无效的文件路径: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.FileInfo().Mode())
			continue
		}

		// 创建文件目录
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// 创建文件
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
	}

	return nil
}

// UploadWorld 上传世界
func (h *WorldHandler) UploadWorld(c *gin.Context) {
	file, header, err := c.Request.FormFile("world")
	if err != nil {
		c.JSON(400, gin.H{"error": "上传文件失败: " + err.Error()})
		return
	}
	defer file.Close()

	// 检查文件扩展名
	filename := header.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".zip") &&
		!strings.HasSuffix(strings.ToLower(filename), ".mcworld") {
		c.JSON(400, gin.H{"error": "只支持 .zip 和 .mcworld 格式"})
		return
	}

	// 获取bedrock路径
	wd, err := os.Getwd()
	if err != nil {
		c.JSON(500, gin.H{"error": "获取工作目录失败: " + err.Error()})
		return
	}
	bedrockPath := filepath.Join(wd, "bedrock-server", "bedrock-server-1.21.95.1")

	// 保存上传的文件
	worldsPath := filepath.Join(bedrockPath, "worlds")
	os.MkdirAll(worldsPath, 0755)

	uploadPath := filepath.Join(worldsPath, filename)
	out, err := os.Create(uploadPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(500, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}

	// 关闭文件句柄以便后续操作
	out.Close()

	// 解压文件
	var extractedWorldName string
	if strings.HasSuffix(strings.ToLower(filename), ".zip") || strings.HasSuffix(strings.ToLower(filename), ".mcworld") {
		// 获取不带扩展名的文件名作为世界名称
		extractedWorldName = strings.TrimSuffix(filename, filepath.Ext(filename))
		extractPath := filepath.Join(worldsPath, extractedWorldName)

		// 解压文件
		if err := extractZip(uploadPath, extractPath); err != nil {
			// 如果解压失败，删除已上传的文件
			os.Remove(uploadPath)
			c.JSON(500, gin.H{"error": "解压文件失败: " + err.Error()})
			return
		}

		// 解压成功后删除原始压缩文件
		if err := os.Remove(uploadPath); err != nil {
			// 记录警告但不影响主流程
			fmt.Printf("警告: 删除压缩文件失败: %v\n", err)
		}

		c.JSON(200, gin.H{"message": fmt.Sprintf("世界文件上传并解压成功: %s", extractedWorldName)})
	} else {
		c.JSON(200, gin.H{"message": "世界文件上传成功: " + filename})
	}
}

// DeleteWorld 删除世界
func (h *WorldHandler) DeleteWorld(c *gin.Context) {
	worldName := c.Param("name")

	if err := h.worldService.DeleteWorld(worldName); err != nil {
		// 根据错误类型返回不同的状态码
		if strings.Contains(err.Error(), "世界不存在") {
			c.JSON(404, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "世界名称不能为空") {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "删除世界失败: " + err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "世界已删除: " + worldName + "，配置文件已同步更新"})
}

// ActivateWorld 激活世界
func (h *WorldHandler) ActivateWorld(c *gin.Context) {
	worldName := c.Param("name")

	if err := h.worldService.ActivateWorld(worldName); err != nil {
		c.JSON(500, gin.H{"error": "激活世界失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "世界已激活: " + worldName + "，重启服务器后生效"})
}