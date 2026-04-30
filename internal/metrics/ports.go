package metrics

import (
	"fmt"
	"net"
	"time"

	"watchup-agent/internal/config"
)

// PortStatus represents the status of a monitored port
type PortStatus struct {
	Port        int           `json:"port"`
	Name        string        `json:"name"`
	Host        string        `json:"host"`
	IsUp        bool          `json:"is_up"`
	ResponseTime time.Duration `json:"response_time_ms"`
	Error       string        `json:"error,omitempty"`
	CheckedAt   time.Time     `json:"checked_at"`
}

// CheckPort checks if a specific port is accessible
func CheckPort(host string, port int, timeout time.Duration) (*PortStatus, error) {
	if host == "" {
		host = "localhost"
	}

	status := &PortStatus{
		Port:      port,
		Host:      host,
		CheckedAt: time.Now(),
	}

	start := time.Now()
	address := fmt.Sprintf("%s:%d", host, port)
	
	conn, err := net.DialTimeout("tcp", address, timeout)
	status.ResponseTime = time.Since(start)

	if err != nil {
		status.IsUp = false
		status.Error = err.Error()
		return status, nil // Not an error - just port is down
	}

	conn.Close()
	status.IsUp = true
	return status, nil
}

// CheckMultiplePorts checks multiple ports concurrently
func CheckMultiplePorts(portChecks []config.PortCheck) ([]*PortStatus, error) {
	if len(portChecks) == 0 {
		return []*PortStatus{}, nil
	}

	results := make([]*PortStatus, len(portChecks))
	resultChan := make(chan struct {
		index  int
		status *PortStatus
		err    error
	}, len(portChecks))

	// Start concurrent checks
	for i, check := range portChecks {
		go func(index int, portCheck config.PortCheck) {
			timeout := 5 * time.Second // Default timeout
			if portCheck.Timeout != "" {
				if parsedTimeout, err := time.ParseDuration(portCheck.Timeout); err == nil {
					timeout = parsedTimeout
				}
			}

			host := portCheck.Host
			if host == "" {
				host = "localhost"
			}

			status, err := CheckPort(host, portCheck.Port, timeout)
			if status != nil {
				status.Name = portCheck.Name
			}

			resultChan <- struct {
				index  int
				status *PortStatus
				err    error
			}{index, status, err}
		}(i, check)
	}

	// Collect results
	for i := 0; i < len(portChecks); i++ {
		result := <-resultChan
		if result.err != nil {
			return nil, fmt.Errorf("failed to check port %d: %w", 
				portChecks[result.index].Port, result.err)
		}
		results[result.index] = result.status
	}

	return results, nil
}

// GetPortSummary returns a summary of port check results
func GetPortSummary(portStatuses []*PortStatus) map[string]interface{} {
	summary := map[string]interface{}{
		"total_ports": len(portStatuses),
		"ports_up":    0,
		"ports_down":  0,
		"avg_response_time_ms": float64(0),
	}

	if len(portStatuses) == 0 {
		return summary
	}

	var totalResponseTime time.Duration
	upCount := 0

	for _, status := range portStatuses {
		if status.IsUp {
			upCount++
			totalResponseTime += status.ResponseTime
		}
	}

	summary["ports_up"] = upCount
	summary["ports_down"] = len(portStatuses) - upCount

	if upCount > 0 {
		avgMs := float64(totalResponseTime.Nanoseconds()) / float64(upCount) / 1000000
		summary["avg_response_time_ms"] = avgMs
	}

	return summary
}