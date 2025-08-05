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
func (p *PerformaceMonitoringHandler) GetPerformaceMonitoring(c *gin.Context) {
	cpu, err := p.performanceMonitoringService.GetCPUUsage()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read configuration: " + err.Error()})
		return
	}
	memory, err := p.performanceMonitoringService.GetMemoryUsage()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read configuration: " + err.Error()})
		return
	}
	disk, err := p.performanceMonitoringService.GetDiskUsage()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read configuration: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"cpu": cpu, "memory": memory, "disk": disk})
}
