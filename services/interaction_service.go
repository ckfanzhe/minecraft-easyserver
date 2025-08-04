package services

import (
	"bufio"
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"

	"minecraft-easyserver/models"
)

type InteractionService struct {
	stdin          io.WriteCloser
	commandHistory []models.ServerCommandResponse
	mutex          sync.RWMutex
	maxHistory     int
	enabled        bool
}

var interactionService *InteractionService

func NewInteractionService() *InteractionService {
	if interactionService == nil {
		interactionService = &InteractionService{
			commandHistory: make([]models.ServerCommandResponse, 0),
			maxHistory:     100, // Keep last 100 command responses
			enabled:        runtime.GOOS != "windows", // Disable on Windows for now
		}
	}
	return interactionService
}

// IsEnabled returns whether interaction is enabled on current platform
func (is *InteractionService) IsEnabled() bool {
	return is.enabled
}

// SetStdin sets the stdin pipe for server interaction
func (is *InteractionService) SetStdin(stdin io.WriteCloser) {
	is.mutex.Lock()
	defer is.mutex.Unlock()
	is.stdin = stdin
}

// SendCommand sends a command to the server
func (is *InteractionService) SendCommand(command string) error {
	if !is.enabled {
		return fmt.Errorf("server interaction is not supported on this platform")
	}

	is.mutex.Lock()
	defer is.mutex.Unlock()

	if is.stdin == nil {
		return fmt.Errorf("server is not running or stdin is not available")
	}

	// Send command to server
	_, err := is.stdin.Write([]byte(command + "\n"))
	if err != nil {
		return fmt.Errorf("failed to send command: %v", err)
	}

	// Record command in history
	response := models.ServerCommandResponse{
		Command:   command,
		Response:  "Command sent", // Will be updated when response is captured
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Success:   true,
	}

	is.addCommandResponse(response)
	return nil
}

// GetCommandHistory returns recent command history
func (is *InteractionService) GetCommandHistory(limit int) []models.ServerCommandResponse {
	is.mutex.RLock()
	defer is.mutex.RUnlock()

	if limit <= 0 || limit > len(is.commandHistory) {
		return is.commandHistory
	}

	start := len(is.commandHistory) - limit
	return is.commandHistory[start:]
}

// ClearHistory clears command history
func (is *InteractionService) ClearHistory() {
	is.mutex.Lock()
	defer is.mutex.Unlock()
	is.commandHistory = make([]models.ServerCommandResponse, 0)
}

// addCommandResponse adds a command response to history
func (is *InteractionService) addCommandResponse(response models.ServerCommandResponse) {
	is.commandHistory = append(is.commandHistory, response)

	// Keep only the last maxHistory entries
	if len(is.commandHistory) > is.maxHistory {
		is.commandHistory = is.commandHistory[len(is.commandHistory)-is.maxHistory:]
	}
}

// Close closes the stdin pipe
func (is *InteractionService) Close() {
	is.mutex.Lock()
	defer is.mutex.Unlock()

	if is.stdin != nil {
		is.stdin.Close()
		is.stdin = nil
	}
}

// ValidateCommand validates if a command is safe to execute
func (is *InteractionService) ValidateCommand(command string) error {
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Add basic validation for dangerous commands
	dangerousCommands := []string{"stop", "restart", "shutdown"}
	for _, dangerous := range dangerousCommands {
		if command == dangerous {
			return fmt.Errorf("command '%s' is not allowed through web interface", command)
		}
	}

	return nil
}

// StartCommandCapture starts capturing command responses from server output
func (is *InteractionService) StartCommandCapture(reader io.ReadCloser) {
	if !is.enabled {
		return
	}

	go is.captureCommandResponses(reader)
}

// captureCommandResponses captures command responses from server output
func (is *InteractionService) captureCommandResponses(reader io.ReadCloser) {
	defer reader.Close()
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		// This is a simplified approach - in a real implementation,
		// you would need to parse the server output to match responses
		// to specific commands based on timing and content patterns
		if line != "" {
			// For now, just log the output
			// In a more sophisticated implementation, you would correlate
			// this output with previously sent commands
		}
	}
}