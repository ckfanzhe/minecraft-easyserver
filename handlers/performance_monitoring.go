package handlers

import (
	"minecraft-easyserver/services"

	"github.com/gin-gonic/gin"
)

// PerformaceMonitoringHandle performance monitoring handler
type PerformaceMonitoringHandler struct {
	performanceMonitoringService *services.PerformanceMonitoringService
}

// PerformaceMonitoringHandler creates a new performance monitoring  handler
func NewPerformaceMonitoringHandler() *PerformaceMonitoringHandler {
	return &PerformaceMonitoringHandler{
		performanceMonitoringService: services.NewPerformanceMonitoringService(),
	}
}

// GetPerformanceMonitoringData gets comprehensive performance monitoring data
func (p *PerformaceMonitoringHandler) GetPerformanceMonitoringData(c *gin.Context) {
	// Get server status to obtain bedrock PID
	serverHandler := NewServerHandler()
	serverStatus := serverHandler.serverService.GetStatus()
	
	bedrockPID := 0
	if serverStatus.Status == "running" {
		bedrockPID = serverStatus.PID
	}

	data, err := p.performanceMonitoringService.GetPerformanceMonitoringData(bedrockPID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get performance monitoring data: " + err.Error()})
		return
	}

	c.JSON(200, data)
}
