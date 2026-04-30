package metrics

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"watchup-agent/internal/config"
)

// LatencyResult represents the result of a latency check
type LatencyResult struct {
	Host         string        `json:"host"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Latency      time.Duration `json:"latency_ms"`
	IsSuccessful bool          `json:"is_successful"`
	Error        string        `json:"error,omitempty"`
	CheckedAt    time.Time     `json:"checked_at"`
}

// CheckTCPLatency checks latency by establishing a TCP connection
func CheckTCPLatency(host string, port int, timeout time.Duration) (*LatencyResult, error) {
	result := &LatencyResult{
		Host:      host,
		Type:      "tcp",
		CheckedAt: time.Now(),
	}

	start := time.Now()
	address := fmt.Sprintf("%s:%d", host, port)
	
	conn, err := net.DialTimeout("tcp", address, timeout)
	result.Latency = time.Since(start)

	if err != nil {
		result.IsSuccessful = false
		result.Error = err.Error()
		return result, nil
	}

	conn.Close()
	result.IsSuccessful = true
	return result, nil
}

// CheckHTTPLatency checks latency by making an HTTP request
func CheckHTTPLatency(url string, timeout time.Duration) (*LatencyResult, error) {
	result := &LatencyResult{
		Host:      url,
		Type:      "http",
		CheckedAt: time.Now(),
	}

	client := &http.Client{
		Timeout: timeout,
	}

	start := time.Now()
	resp, err := client.Get(url)
	result.Latency = time.Since(start)

	if err != nil {
		result.IsSuccessful = false
		result.Error = err.Error()
		return result, nil
	}

	resp.Body.Close()
	result.IsSuccessful = resp.StatusCode >= 200 && resp.StatusCode < 400
	
	if !result.IsSuccessful {
		result.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	return result, nil
}

// CheckMultipleLatencies checks multiple latency targets concurrently
func CheckMultipleLatencies(latencyChecks []config.LatencyCheck) ([]*LatencyResult, error) {
	if len(latencyChecks) == 0 {
		return []*LatencyResult{}, nil
	}

	results := make([]*LatencyResult, len(latencyChecks))
	resultChan := make(chan struct {
		index  int
		result *LatencyResult
		err    error
	}, len(latencyChecks))

	// Start concurrent checks
	for i, check := range latencyChecks {
		go func(index int, latencyCheck config.LatencyCheck) {
			timeout := 5 * time.Second // Default timeout
			if latencyCheck.Timeout != "" {
				if parsedTimeout, err := time.ParseDuration(latencyCheck.Timeout); err == nil {
					timeout = parsedTimeout
				}
			}

			var result *LatencyResult
			var err error

			switch latencyCheck.Type {
			case "http":
				url := latencyCheck.URL
				if url == "" {
					url = fmt.Sprintf("http://%s", latencyCheck.Host)
				}
				result, err = CheckHTTPLatency(url, timeout)
			case "tcp":
				port := latencyCheck.Port
				if port == 0 {
					port = 80 // Default port
				}
				result, err = CheckTCPLatency(latencyCheck.Host, port, timeout)
			default:
				// Default to TCP check
				port := latencyCheck.Port
				if port == 0 {
					port = 80
				}
				result, err = CheckTCPLatency(latencyCheck.Host, port, timeout)
			}

			if result != nil {
				result.Name = latencyCheck.Name
			}

			resultChan <- struct {
				index  int
				result *LatencyResult
				err    error
			}{index, result, err}
		}(i, check)
	}

	// Collect results
	for i := 0; i < len(latencyChecks); i++ {
		res := <-resultChan
		if res.err != nil {
			return nil, fmt.Errorf("failed to check latency for %s: %w", 
				latencyChecks[res.index].Host, res.err)
		}
		results[res.index] = res.result
	}

	return results, nil
}

// GetLatencySummary returns a summary of latency check results
func GetLatencySummary(latencyResults []*LatencyResult) map[string]interface{} {
	summary := map[string]interface{}{
		"total_checks":     len(latencyResults),
		"successful":       0,
		"failed":          0,
		"avg_latency_ms":  float64(0),
		"min_latency_ms":  float64(0),
		"max_latency_ms":  float64(0),
	}

	if len(latencyResults) == 0 {
		return summary
	}

	var totalLatency time.Duration
	var minLatency, maxLatency time.Duration
	successCount := 0
	first := true

	for _, result := range latencyResults {
		if result.IsSuccessful {
			successCount++
			totalLatency += result.Latency

			if first || result.Latency < minLatency {
				minLatency = result.Latency
			}
			if first || result.Latency > maxLatency {
				maxLatency = result.Latency
			}
			first = false
		}
	}

	summary["successful"] = successCount
	summary["failed"] = len(latencyResults) - successCount

	if successCount > 0 {
		avgMs := float64(totalLatency.Nanoseconds()) / float64(successCount) / 1000000
		summary["avg_latency_ms"] = avgMs
		summary["min_latency_ms"] = float64(minLatency.Nanoseconds()) / 1000000
		summary["max_latency_ms"] = float64(maxLatency.Nanoseconds()) / 1000000
	}

	return summary
}