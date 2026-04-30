package metrics

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/host"
)

// SystemMetrics represents general system information
type SystemMetrics struct {
	Uptime          uint64 `json:"uptime"`           // System uptime in seconds
	Hostname        string `json:"hostname"`         // System hostname
	OS              string `json:"os"`               // Operating system
	Platform        string `json:"platform"`         // Platform (e.g., ubuntu, windows)
	PlatformVersion string `json:"platform_version"` // Platform version
	Architecture    string `json:"architecture"`     // System architecture
	BootTime        uint64 `json:"boot_time"`        // Boot time (Unix timestamp)
}

// GetSystemInfo collects general system information
func GetSystemInfo() (*SystemMetrics, error) {
	// Get host information
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %w", err)
	}

	return &SystemMetrics{
		Uptime:          hostInfo.Uptime,
		Hostname:        hostInfo.Hostname,
		OS:              hostInfo.OS,
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		Architecture:    runtime.GOARCH,
		BootTime:        hostInfo.BootTime,
	}, nil
}

// GetUptime returns just the system uptime in seconds
func GetUptime() (uint64, error) {
	uptime, err := host.Uptime()
	if err != nil {
		return 0, fmt.Errorf("failed to get uptime: %w", err)
	}
	return uptime, nil
}