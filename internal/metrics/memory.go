package metrics

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
)

// MemoryMetrics represents memory usage information
type MemoryMetrics struct {
	Total       uint64  `json:"total"`        // Total memory in bytes
	Used        uint64  `json:"used"`         // Used memory in bytes
	Available   uint64  `json:"available"`    // Available memory in bytes
	UsedPercent float64 `json:"used_percent"` // Used percentage
	Free        uint64  `json:"free"`         // Free memory in bytes
	Cached      uint64  `json:"cached"`       // Cached memory in bytes (Linux)
	Buffers     uint64  `json:"buffers"`      // Buffer memory in bytes (Linux)
}

// GetMemoryUsage collects memory usage metrics
func GetMemoryUsage() (*MemoryMetrics, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory usage: %w", err)
	}

	return &MemoryMetrics{
		Total:       vmStat.Total,
		Used:        vmStat.Used,
		Available:   vmStat.Available,
		UsedPercent: vmStat.UsedPercent,
		Free:        vmStat.Free,
		Cached:      vmStat.Cached,
		Buffers:     vmStat.Buffers,
	}, nil
}

// GetSwapUsage collects swap memory usage (optional)
func GetSwapUsage() (map[string]interface{}, error) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		// Swap might not be available on all systems
		return map[string]interface{}{
			"total":        uint64(0),
			"used":         uint64(0),
			"free":         uint64(0),
			"used_percent": float64(0),
		}, nil
	}

	return map[string]interface{}{
		"total":        swapStat.Total,
		"used":         swapStat.Used,
		"free":         swapStat.Free,
		"used_percent": swapStat.UsedPercent,
	}, nil
}