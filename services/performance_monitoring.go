package services

import (
	"fmt"
	"minecraft-easyserver/models"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// PerformanceMonitoringService performance monitoring
type PerformanceMonitoringService struct{}

// NewPerformanceMonitoringService creates a new performance monitoring instance
func NewPerformanceMonitoringService() *PerformanceMonitoringService {
	return &PerformanceMonitoringService{}
}

// GetCPUUsage gets system CPU usage
func (p *PerformanceMonitoringService) GetCPUUsage() (float64, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Printf("Failed to obtain CPU usage rate: %v\n", err)
		return 0, err
	}
	return percent[0], nil
}

// GetMemoryUsage gets system memory usage
func (p *PerformanceMonitoringService) GetMemoryUsage() (float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("Failed to retrieve memory information: %v\n", err)
		return 0, err
	}
	return v.UsedPercent, nil
}

// GetDiskUsage gets disk usage (kept for backward compatibility)
func (p *PerformanceMonitoringService) GetDiskUsage() (float64, error) {
	// Get current executable file disk
	exePath, err := os.Executable()
	if err != nil {
		return 0, fmt.Errorf("Failed to obtain executable file path: %v", err)
	}

	var mountPoint string
	if runtime.GOOS == "windows" {
		mountPoint = filepath.VolumeName(exePath)
		if mountPoint == "" {
			mountPoint = "C:\\"
		}
	} else {
		mountPoint = "/"
	}

	usage, err := disk.Usage(mountPoint)
	if err != nil {
		return 0, fmt.Errorf("Failed to retrieve disk usage status: %v", err)
	}
	return usage.UsedPercent, nil
}

// GetSystemPerformance gets system performance data
func (p *PerformanceMonitoringService) GetSystemPerformance() (models.SystemPerformance, error) {
	cpuUsage, err := p.GetCPUUsage()
	if err != nil {
		return models.SystemPerformance{}, err
	}

	memoryUsage, err := p.GetMemoryUsage()
	if err != nil {
		return models.SystemPerformance{}, err
	}

	return models.SystemPerformance{
		CPUUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
		Timestamp:   time.Now().Format(time.RFC3339),
	}, nil
}

// GetBedrockProcessPerformance gets bedrock process performance data
func (p *PerformanceMonitoringService) GetBedrockProcessPerformance(pid int) (models.ProcessPerformance, error) {
	if pid <= 0 {
		return models.ProcessPerformance{
			PID:       0,
			CPUUsage:  0,
			MemoryUsage: 0,
			MemoryMB:  0,
			Timestamp: time.Now().Format(time.RFC3339),
		}, nil
	}

	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return models.ProcessPerformance{}, fmt.Errorf("Failed to get process info: %v", err)
	}

	// Get CPU usage
	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		cpuPercent = 0
	}

	// Get memory info
	memInfo, err := proc.MemoryInfo()
	if err != nil {
		return models.ProcessPerformance{}, fmt.Errorf("Failed to get memory info: %v", err)
	}

	// Get system memory to calculate percentage
	sysMem, err := mem.VirtualMemory()
	if err != nil {
		return models.ProcessPerformance{}, fmt.Errorf("Failed to get system memory: %v", err)
	}

	memoryUsagePercent := float64(memInfo.RSS) / float64(sysMem.Total) * 100
	memoryMB := float64(memInfo.RSS) / 1024 / 1024

	return models.ProcessPerformance{
		PID:         pid,
		CPUUsage:    cpuPercent,
		MemoryUsage: memoryUsagePercent,
		MemoryMB:    memoryMB,
		Timestamp:   time.Now().Format(time.RFC3339),
	}, nil
}

// GetPerformanceMonitoringData gets combined performance monitoring data
func (p *PerformanceMonitoringService) GetPerformanceMonitoringData(bedrockPID int) (models.PerformanceMonitoringData, error) {
	systemPerf, err := p.GetSystemPerformance()
	if err != nil {
		return models.PerformanceMonitoringData{}, err
	}

	bedrockPerf, err := p.GetBedrockProcessPerformance(bedrockPID)
	if err != nil {
		// If bedrock process monitoring fails, return empty process data
		bedrockPerf = models.ProcessPerformance{
			PID:       0,
			CPUUsage:  0,
			MemoryUsage: 0,
			MemoryMB:  0,
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	return models.PerformanceMonitoringData{
		System:  systemPerf,
		Bedrock: bedrockPerf,
	}, nil
}
