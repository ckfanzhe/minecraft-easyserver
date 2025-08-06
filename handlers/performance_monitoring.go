package handlers

import (
	"minecraft-easyserver/services"

	"github.com/gin-gonic/gin"
)

// PerformanceMonitoringHandler performance monitoring handler
type PerformanceMonitoringHandler struct {
	performanceMonitoringService *services.PerformanceMonitoringService
}

// NewPerformanceMonitoringHandler creates a new performance monitoring handler
func NewPerformanceMonitoringHandler() *PerformanceMonitoringHandler {
	return &PerformanceMonitoringHandler{
		performanceMonitoringService: services.NewPerformanceMonitoringService(),
	}
}

// GetPerformanceMonitoringData gets comprehensive performance monitoring data
func (p *PerformanceMonitoringHandler) GetPerformanceMonitoringData(c *gin.Context) {
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
