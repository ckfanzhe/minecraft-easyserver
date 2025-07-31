package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config application configuration structure
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
	// AppConfig global configuration instance
	AppConfig *Config
)

// LoadConfig loads configuration from file
func LoadConfig(configPath string) error {
	// If configuration file doesn't exist, create default configuration file
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath); err != nil {
			return fmt.Errorf("failed to create default configuration file: %v", err)
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %v", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse configuration file: %v", err)
	}

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return fmt.Errorf("configuration validation failed: %v", err)
	}

	AppConfig = config
	return nil
}

// createDefaultConfig creates default configuration file
func createDefaultConfig(configPath string) error {
	defaultConfig := &Config{}
	
	// Set default values
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

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// validateConfig validates configuration
func validateConfig(config *Config) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Bedrock.Path == "" {
		return fmt.Errorf("Bedrock path cannot be empty")
	}

	if config.Bedrock.Executable == "" {
		return fmt.Errorf("Bedrock executable name cannot be empty")
	}

	return nil
}

// GetServerAddress gets server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetBedrockExecutablePath gets Bedrock executable full path
func (c *Config) GetBedrockExecutablePath() string {
	return filepath.Join(c.Bedrock.Path, c.Bedrock.Executable)
}

// GetBedrockPath gets Bedrock path
func (c *Config) GetBedrockPath() string {
	return c.Bedrock.Path
}