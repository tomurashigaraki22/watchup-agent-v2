package metrics

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// CPUMetrics represents CPU usage information
type CPUMetrics struct {
	UsagePercent float64   `json:"usage_percent"`
	PerCore      []float64 `json:"per_core,omitempty"`
	Cores        int       `json:"cores"`
}

// GetCPUUsage collects CPU usage metrics
func GetCPUUsage() (*CPUMetrics, error) {
	// Get overall CPU usage (1 second sample)
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %w", err)
	}

	if len(percentages) == 0 {
		return nil, fmt.Errorf("no CPU usage data available")
	}

	// Get per-core usage
	perCore, err := cpu.Percent(0, true) // Use cached data for per-core
	if err != nil {
		// Per-core data is optional, continue without it
		perCore = nil
	}

	// Get CPU count
	coreCount, err := cpu.Counts(true) // Logical cores
	if err != nil {
		coreCount = len(perCore) // Fallback to per-core count
	}

	return &CPUMetrics{
		UsagePercent: percentages[0],
		PerCore:      perCore,
		Cores:        coreCount,
	}, nil
}

// GetCPUInfo returns basic CPU information (for initial setup)
func GetCPUInfo() (map[string]interface{}, error) {
	info, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %w", err)
	}

	if len(info) == 0 {
		return map[string]interface{}{
			"model": "Unknown",
			"cores": 0,
		}, nil
	}

	return map[string]interface{}{
		"model":     info[0].ModelName,
		"cores":     info[0].Cores,
		"mhz":       info[0].Mhz,
		"vendor_id": info[0].VendorID,
	}, nil
}