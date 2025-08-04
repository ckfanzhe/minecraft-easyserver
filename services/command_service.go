package services

import (
	"fmt"
	"sync"

	"minecraft-easyserver/models"
)

type CommandService struct {
	quickCommands []models.QuickCommand
	mutex         sync.RWMutex
}

var commandService *CommandService

func NewCommandService() *CommandService {
	if commandService == nil {
		commandService = &CommandService{
			quickCommands: initializeQuickCommands(),
		}
	}
	return commandService
}

// initializeQuickCommands initializes the default quick commands
func initializeQuickCommands() []models.QuickCommand {
	return []models.QuickCommand{
		// Time commands
		{
			ID:          "time_day",
			Name:        "设置为白天",
			Description: "将游戏时间设置为白天",
			Command:     "time set day",
			Category:    "time",
		},
		{
			ID:          "time_night",
			Name:        "设置为夜晚",
			Description: "将游戏时间设置为夜晚",
			Command:     "time set night",
			Category:    "time",
		},
		{
			ID:          "time_noon",
			Name:        "设置为正午",
			Description: "将游戏时间设置为正午",
			Command:     "time set noon",
			Category:    "time",
		},
		{
			ID:          "time_midnight",
			Name:        "设置为午夜",
			Description: "将游戏时间设置为午夜",
			Command:     "time set midnight",
			Category:    "time",
		},
		// Weather commands
		{
			ID:          "weather_clear",
			Name:        "晴天",
			Description: "设置天气为晴天",
			Command:     "weather clear",
			Category:    "weather",
		},
		{
			ID:          "weather_rain",
			Name:        "雨天",
			Description: "设置天气为雨天",
			Command:     "weather rain",
			Category:    "weather",
		},
		{
			ID:          "weather_thunder",
			Name:        "雷阵雨",
			Description: "设置天气为雷阵雨",
			Command:     "weather thunder",
			Category:    "weather",
		},
		// Game mode commands
		{
			ID:          "gamemode_survival",
			Name:        "生存模式",
			Description: "将默认游戏模式设置为生存模式",
			Command:     "gamemode survival",
			Category:    "gamemode",
		},
		{
			ID:          "gamemode_creative",
			Name:        "创造模式",
			Description: "将默认游戏模式设置为创造模式",
			Command:     "gamemode creative",
			Category:    "gamemode",
		},
		{
			ID:          "gamemode_adventure",
			Name:        "冒险模式",
			Description: "将默认游戏模式设置为冒险模式",
			Command:     "gamemode adventure",
			Category:    "gamemode",
		},
		// Difficulty commands
		{
			ID:          "difficulty_peaceful",
			Name:        "和平难度",
			Description: "设置游戏难度为和平",
			Command:     "difficulty peaceful",
			Category:    "difficulty",
		},
		{
			ID:          "difficulty_easy",
			Name:        "简单难度",
			Description: "设置游戏难度为简单",
			Command:     "difficulty easy",
			Category:    "difficulty",
		},
		{
			ID:          "difficulty_normal",
			Name:        "普通难度",
			Description: "设置游戏难度为普通",
			Command:     "difficulty normal",
			Category:    "difficulty",
		},
		{
			ID:          "difficulty_hard",
			Name:        "困难难度",
			Description: "设置游戏难度为困难",
			Command:     "difficulty hard",
			Category:    "difficulty",
		},
	}
}

// GetQuickCommands returns all quick commands
func (cs *CommandService) GetQuickCommands() []models.QuickCommand {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	return cs.quickCommands
}

// GetQuickCommandsByCategory returns quick commands filtered by category
func (cs *CommandService) GetQuickCommandsByCategory(category string) []models.QuickCommand {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	var filtered []models.QuickCommand
	for _, cmd := range cs.quickCommands {
		if cmd.Category == category {
			filtered = append(filtered, cmd)
		}
	}
	return filtered
}

// GetQuickCommandByID returns a quick command by ID
func (cs *CommandService) GetQuickCommandByID(id string) (*models.QuickCommand, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	for _, cmd := range cs.quickCommands {
		if cmd.ID == id {
			return &cmd, nil
		}
	}
	return nil, fmt.Errorf("quick command with ID '%s' not found", id)
}

// ExecuteQuickCommand executes a quick command by ID
func (cs *CommandService) ExecuteQuickCommand(id string) error {
	cmd, err := cs.GetQuickCommandByID(id)
	if err != nil {
		return err
	}

	// Get interaction service to send the command
	interactionSvc := NewInteractionService()
	if !interactionSvc.IsEnabled() {
		return fmt.Errorf("server interaction is not supported on this platform")
	}

	// Validate and send the command
	if err := interactionSvc.ValidateCommand(cmd.Command); err != nil {
		return fmt.Errorf("command validation failed: %v", err)
	}

	return interactionSvc.SendCommand(cmd.Command)
}

// AddQuickCommand adds a new quick command
func (cs *CommandService) AddQuickCommand(cmd models.QuickCommand) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	// Check if ID already exists
	for _, existing := range cs.quickCommands {
		if existing.ID == cmd.ID {
			return fmt.Errorf("quick command with ID '%s' already exists", cmd.ID)
		}
	}

	cs.quickCommands = append(cs.quickCommands, cmd)
	return nil
}

// RemoveQuickCommand removes a quick command by ID
func (cs *CommandService) RemoveQuickCommand(id string) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for i, cmd := range cs.quickCommands {
		if cmd.ID == id {
			cs.quickCommands = append(cs.quickCommands[:i], cs.quickCommands[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("quick command with ID '%s' not found", id)
}

// GetCategories returns all available command categories
func (cs *CommandService) GetCategories() []string {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	categoryMap := make(map[string]bool)
	for _, cmd := range cs.quickCommands {
		categoryMap[cmd.Category] = true
	}

	var categories []string
	for category := range categoryMap {
		categories = append(categories, category)
	}
	return categories
}