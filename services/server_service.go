package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"minecraft-easyserver/models"
)

var (
	serverProcess *exec.Cmd
	serverMutex   sync.Mutex
	bedrockPath   string
	logSvc        *LogService
	interactionSvc *InteractionService
)

// InitBedrockPath initializes bedrock path
func InitBedrockPath(path string) error {
	if path == "" {
		return fmt.Errorf("bedrock path cannot be empty")
	}

	// If it's a relative path, convert to absolute path
	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		path = filepath.Join(wd, path)
	}

	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("bedrock path does not exist: %s", path)
	}

	bedrockPath = path
	return nil
}

// SetBedrockPath sets bedrock path (mainly for testing)
func SetBedrockPath(path string) {
	bedrockPath = path
}

// GetBedrockPath returns the current bedrock path
func GetBedrockPath() string {
	return bedrockPath
}

// ServerService server service
type ServerService struct{}

// NewServerService creates a new server service instance
func NewServerService() *ServerService {
	return &ServerService{}
}

// GetStatus gets server status
func (s *ServerService) GetStatus() models.ServerStatus {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverProcess == nil || serverProcess.Process == nil {
		return models.ServerStatus{
			Status:  "stopped",
			Message: "Server not running",
		}
	}

	// Check if process is still running
	process, err := os.FindProcess(serverProcess.Process.Pid)
	if err != nil {
		serverProcess = nil
		return models.ServerStatus{
			Status:  "stopped",
			Message: "Server not running",
		}
	}

	// On Windows, simply check if process exists
	// If process has ended, FindProcess will still return a Process object
	// We can try sending signal 0 to check if process is really running
	if process != nil {
		return models.ServerStatus{
			Status:  "running",
			Message: "Server is running",
			PID:     serverProcess.Process.Pid,
		}
	}

	serverProcess = nil
	return models.ServerStatus{
		Status:  "stopped",
		Message: "Server not running",
	}
}

// Start starts server
func (s *ServerService) Start() error {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverProcess != nil && serverProcess.Process != nil {
		return fmt.Errorf("server is already running")
	}

	// Check if bedrock path is configured
	if bedrockPath == "" {
		return fmt.Errorf("no server version is currently active. Please download and activate a server version first")
	}

	// Initialize services
	logSvc = NewLogService()
	interactionSvc = NewInteractionService()

	// Get executable name based on operating system
	var executableName string
	if runtime.GOOS == "windows" {
		executableName = "bedrock_server.exe"
	} else {
		executableName = "bedrock_server"
	}
	
	exePath := filepath.Join(bedrockPath, executableName)
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		return fmt.Errorf("%s file not found in %s. Please ensure the server version is properly downloaded", executableName, bedrockPath)
	}

	serverProcess = exec.Command(exePath)
	serverProcess.Dir = bedrockPath

	// Set up pipes for logging and interaction
	stdout, err := serverProcess.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := serverProcess.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Set up stdin for interaction (only on supported platforms)
	if interactionSvc.IsEnabled() {
		stdin, err := serverProcess.StdinPipe()
		if err != nil {
			return fmt.Errorf("failed to create stdin pipe: %v", err)
		}
		interactionSvc.SetStdin(stdin)
	}

	if err := serverProcess.Start(); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	// Start log capture
	logSvc.StartLogCapture(stdout, stderr)
	logSvc.AddLogEntry("INFO", "Server started successfully")

	// Start command response capture (if interaction is enabled)
	if interactionSvc.IsEnabled() {
		// Note: In a real implementation, you might want to duplicate stdout
		// to capture both logs and command responses separately
		logSvc.AddLogEntry("INFO", "Server interaction enabled")
	}

	return nil
}

// Stop stops server
func (s *ServerService) Stop() error {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverProcess == nil || serverProcess.Process == nil {
		return fmt.Errorf("server not running")
	}

	// Log server stopping
	if logSvc != nil {
		logSvc.AddLogEntry("INFO", "Stopping server...")
	}

	// Stop log capture
	if logSvc != nil {
		logSvc.StopLogCapture()
	}

	// Close interaction service
	if interactionSvc != nil {
		interactionSvc.Close()
	}

	if err := serverProcess.Process.Kill(); err != nil {
		return fmt.Errorf("failed to stop server: %v", err)
	}

	serverProcess.Wait()
	serverProcess = nil

	// Log server stopped
	if logSvc != nil {
		logSvc.AddLogEntry("INFO", "Server stopped")
	}

	return nil
}

// Restart restarts server
func (s *ServerService) Restart() error {
	// Stop first
	err := s.Stop()
	if err != nil {
		// If stop fails, we still try to start
		// Only log if logSvc is available
		if logSvc != nil {
			logSvc.AddLogEntry("WARN", fmt.Sprintf("Failed to stop server gracefully: %v", err))
		}
	}

	// Wait one second
	time.Sleep(time.Second)

	// Start again (this will reinitialize all services)
	return s.Start()
}