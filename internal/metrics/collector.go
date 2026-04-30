package metrics

import (
	"fmt"
	"log"
	"time"

	"watchup-agent/internal/config"
)

// MetricsPayload represents the complete metrics payload sent to the server
type MetricsPayload struct {
	ServerID  string                 `json:"server_id"`
	Timestamp int64                  `json:"timestamp"`
	Metrics   map[string]interface{} `json:"metrics"`
}

// Collector manages metrics collection
type Collector struct {
	config *config.Config
}

// NewCollector creates a new metrics collector
func NewCollector(cfg *config.Config) *Collector {
	return &Collector{
		config: cfg,
	}
}

// CollectAll collects all enabled metrics and returns a payload
func (c *Collector) CollectAll() (*MetricsPayload, error) {
	metrics := make(map[string]interface{})
	
	// Collect CPU metrics
	if c.config.Metrics.CPU {
		if cpuMetrics, err := GetCPUUsage(); err != nil {
			log.Printf("Warning: Failed to collect CPU metrics: %v", err)
		} else {
			metrics["cpu"] = cpuMetrics
		}
	}

	// Collect Memory metrics
	if c.config.Metrics.Memory {
		if memMetrics, err := GetMemoryUsage(); err != nil {
			log.Printf("Warning: Failed to collect memory metrics: %v", err)
		} else {
			metrics["memory"] = memMetrics
		}

		// Add swap information
		if swapMetrics, err := GetSwapUsage(); err != nil {
			log.Printf("Warning: Failed to collect swap metrics: %v", err)
		} else {
			metrics["swap"] = swapMetrics
		}
	}

	// Collect Disk metrics
	if c.config.Metrics.Disk {
		if diskMetrics, err := GetDiskUsage(); err != nil {
			log.Printf("Warning: Failed to collect disk metrics: %v", err)
		} else {
			metrics["disk"] = diskMetrics
		}

		// Add disk I/O information
		if diskIO, err := GetDiskIO(); err != nil {
			log.Printf("Warning: Failed to collect disk I/O metrics: %v", err)
		} else {
			metrics["disk_io"] = diskIO
		}
	}

	// Collect Network metrics
	if c.config.Metrics.Network {
		if networkMetrics, err := GetNetworkUsage(); err != nil {
			log.Printf("Warning: Failed to collect network metrics: %v", err)
		} else {
			metrics["network_interfaces"] = networkMetrics
		}

		// Add network summary
		if networkSummary, err := GetNetworkSummary(); err != nil {
			log.Printf("Warning: Failed to collect network summary: %v", err)
		} else {
			metrics["network_summary"] = networkSummary
		}
	}

	// Collect Connection metrics (Phase 4)
	if c.config.Metrics.Connections {
		if connMetrics, err := GetActiveConnections(); err != nil {
			log.Printf("Warning: Failed to collect connection metrics: %v", err)
		} else {
			metrics["connections"] = connMetrics
		}

		// Add listening ports
		if listeningPorts, err := GetListeningPorts(); err != nil {
			log.Printf("Warning: Failed to collect listening ports: %v", err)
		} else {
			metrics["listening_ports"] = listeningPorts
		}
	}

	// Collect Port monitoring (Phase 4)
	if len(c.config.Ports) > 0 {
		if portStatuses, err := CheckMultiplePorts(c.config.Ports); err != nil {
			log.Printf("Warning: Failed to check ports: %v", err)
		} else {
			metrics["port_checks"] = portStatuses
			metrics["port_summary"] = GetPortSummary(portStatuses)
		}
	}

	// Collect Latency checks (Phase 4)
	if len(c.config.LatencyChecks) > 0 {
		if latencyResults, err := CheckMultipleLatencies(c.config.LatencyChecks); err != nil {
			log.Printf("Warning: Failed to check latencies: %v", err)
		} else {
			metrics["latency_checks"] = latencyResults
			metrics["latency_summary"] = GetLatencySummary(latencyResults)
		}
	}

	// Always collect system information
	if systemMetrics, err := GetSystemInfo(); err != nil {
		log.Printf("Warning: Failed to collect system metrics: %v", err)
	} else {
		metrics["system"] = systemMetrics
	}

	// Check if we collected any metrics
	if len(metrics) == 0 {
		return nil, fmt.Errorf("no metrics were collected successfully")
	}

	return &MetricsPayload{
		ServerID:  c.config.ServerID,
		Timestamp: time.Now().Unix(),
		Metrics:   metrics,
	}, nil
}

// CollectSystemInfo collects one-time system information (for registration)
func CollectSystemInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	// System info
	if systemInfo, err := GetSystemInfo(); err == nil {
		info["system"] = systemInfo
	}

	// CPU info
	if cpuInfo, err := GetCPUInfo(); err == nil {
		info["cpu_info"] = cpuInfo
	}

	return info, nil
}