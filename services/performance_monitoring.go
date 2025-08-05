package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// PerformanceMonitoringService performance monitoring
type PerformanceMonitoringService struct{}

// NewPerformanceMonitoringService creates a new performance monitoring instance
func NewPerformanceMonitoringService() *PerformanceMonitoringService {
	return &PerformanceMonitoringService{}
}

func (p *PerformanceMonitoringService) GetCPUUsage() (float64, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Printf("Failed to obtain CPU usage rate: %v\n", err)
		return 0, err
	}
	return percent[0], nil
}

func (p *PerformanceMonitoringService) GetMemoryUsage() (float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("Failed to retrieve memory information: %v\n", err)
		return 0, err
	}
	return v.UsedPercent, nil
}

func (p *PerformanceMonitoringService) GetDiskUsage() (float64, error) {
	// 获取当前可执行文件所在的磁盘
	exePath, err := os.Executable()
	if err != nil {
		return 0, fmt.Errorf("Failed to obtain executable file path: %v", err)
	}

	mountPoint := filepath.VolumeName(exePath)
	if mountPoint == "" {
		return 0, fmt.Errorf("Unable to determine the current disk mount point")
	}

	usage, err := disk.Usage(mountPoint)
	if err != nil {
		return 0, fmt.Errorf("Failed to retrieve disk usage status: %v", err)
	}
	return usage.UsedPercent, nil
}
