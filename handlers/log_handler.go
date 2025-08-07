package handlers

import (
	"log"
	"net/http"
	"strconv"

	"minecraft-easyserver/services"
	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	logService *services.LogService
}

func NewLogHandler() *LogHandler {
	return &LogHandler{
		logService: services.NewLogService(),
	}
}

// GetLogs handles GET /api/logs
func (h *LogHandler) GetLogs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	logs := h.logService.GetLogs(limit)
	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"count": len(logs),
	})
}

// ClearLogs handles DELETE /api/logs
func (h *LogHandler) ClearLogs(c *gin.Context) {
	h.logService.ClearLogs()
	c.JSON(http.StatusOK, gin.H{
		"message": "Logs cleared successfully",
	})
}

// HandleWebSocket handles WebSocket connections for real-time logs
func (h *LogHandler) HandleWebSocket(c *gin.Context) {
	// Token validation is already done by middleware
	// Just proceed with WebSocket upgrade
	h.logService.HandleWebSocket(c.Writer, c.Request)
}

// HandleWebSocketWithAuth handles WebSocket connections with JWT authentication
func (h *LogHandler) HandleWebSocketWithAuth(c *gin.Context) {
	log.Printf("=== HandleWebSocketWithAuth called ===")
	// Get token from query parameter
	tokenString := c.Query("token")
	log.Printf("WebSocket auth: received token: %s", tokenString[:20]+"...")
	if tokenString == "" {
		log.Printf("WebSocket auth: no token provided")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token required in query parameter",
		})
		return
	}

	// Validate the token using auth service
	authService := services.NewAuthService()
	claims, err := authService.ValidateJWT(tokenString)
	if err != nil {
		log.Printf("WebSocket auth: token validation failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired token: " + err.Error(),
		})
		return
	}

	log.Printf("WebSocket auth: token validated successfully for claims: %+v", claims)
	// If token is valid, proceed with WebSocket upgrade
	h.logService.HandleWebSocket(c.Writer, c.Request)
}