package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 应用程序配置结构
type Config struct {
	Server struct {
		Host  string `yaml:"host"`
		Port  int    `yaml:"port"`
		Debug bool   `yaml:"debug"`
	} `yaml:"server"`

	Bedrock struct {
		Path       string `yaml:"path"`
		Executable string `yaml:"executable"`
	} `yaml:"bedrock"`

	Web struct {
		StaticDir    string `yaml:"static_dir"`
		TemplateFile string `yaml:"template_file"`
	} `yaml:"web"`

	Logging struct {
		Level      string `yaml:"level"`
		FileOutput bool   `yaml:"file_output"`
		FilePath   string `yaml:"file_path"`
	} `yaml:"logging"`
}

var (
	// AppConfig 全局配置实例
	AppConfig *Config
)

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) error {
	// 如果配置文件不存在，创建默认配置文件
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath); err != nil {
			return fmt.Errorf("创建默认配置文件失败: %v", err)
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置
	if err := validateConfig(config); err != nil {
		return fmt.Errorf("配置验证失败: %v", err)
	}

	AppConfig = config
	return nil
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig(configPath string) error {
	defaultConfig := &Config{}
	
	// 设置默认值
	defaultConfig.Server.Host = "localhost"
	defaultConfig.Server.Port = 8080
	defaultConfig.Server.Debug = false
	
	defaultConfig.Bedrock.Path = "./bedrock-server/bedrock-server-1.21.95.1"
	defaultConfig.Bedrock.Executable = "bedrock_server.exe"
	
	defaultConfig.Web.StaticDir = "./web"
	defaultConfig.Web.TemplateFile = "./web/index.html"
	
	defaultConfig.Logging.Level = "info"
	defaultConfig.Logging.FileOutput = false
	defaultConfig.Logging.FilePath = "./logs/server.log"

	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("无效的服务器端口: %d", config.Server.Port)
	}

	if config.Bedrock.Path == "" {
		return fmt.Errorf("Bedrock路径不能为空")
	}

	if config.Bedrock.Executable == "" {
		return fmt.Errorf("Bedrock可执行文件名不能为空")
	}

	return nil
}

// GetServerAddress 获取服务器地址
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetBedrockExecutablePath 获取Bedrock可执行文件完整路径
func (c *Config) GetBedrockExecutablePath() string {
	return filepath.Join(c.Bedrock.Path, c.Bedrock.Executable)
}

// GetBedrockPath 获取Bedrock路径
func (c *Config) GetBedrockPath() string {
	return c.Bedrock.Path
}