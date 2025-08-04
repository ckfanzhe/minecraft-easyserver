package handlers

import (
	"net/http"

	"minecraft-easyserver/models"
	"minecraft-easyserver/services"
	"github.com/gin-gonic/gin"
)

type CommandHandler struct {
	commandService *services.CommandService
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		commandService: services.NewCommandService(),
	}
}

// GetQuickCommands handles GET /api/commands
func (h *CommandHandler) GetQuickCommands(c *gin.Context) {
	category := c.Query("category")

	var commands []models.QuickCommand
	if category != "" {
		commands = h.commandService.GetQuickCommandsByCategory(category)
	} else {
		commands = h.commandService.GetQuickCommands()
	}

	c.JSON(http.StatusOK, gin.H{
		"commands": commands,
		"count": len(commands),
	})
}

// GetCategories handles GET /api/commands/categories
func (h *CommandHandler) GetCategories(c *gin.Context) {
	categories := h.commandService.GetCategories()
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"count": len(categories),
	})
}

// ExecuteQuickCommand handles POST /api/commands/:id/execute
func (h *CommandHandler) ExecuteQuickCommand(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Command ID is required",
		})
		return
	}

	// Get command details
	cmd, err := h.commandService.GetQuickCommandByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Execute command
	if err := h.commandService.ExecuteQuickCommand(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Command executed successfully",
		"command": cmd.Command,
		"name": cmd.Name,
	})
}

// AddQuickCommand handles POST /api/commands
func (h *CommandHandler) AddQuickCommand(c *gin.Context) {
	var cmd models.QuickCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if err := h.commandService.AddQuickCommand(cmd); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Quick command added successfully",
		"command": cmd,
	})
}

// RemoveQuickCommand handles DELETE /api/commands/:id
func (h *CommandHandler) RemoveQuickCommand(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Command ID is required",
		})
		return
	}

	if err := h.commandService.RemoveQuickCommand(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quick command removed successfully",
	})
}