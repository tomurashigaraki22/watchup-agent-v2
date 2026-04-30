package metrics

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/net"
)

// ConnectionMetrics represents network connection statistics
type ConnectionMetrics struct {
	Total       int            `json:"total"`        // Total connections
	TCP         int            `json:"tcp"`          // TCP connections
	UDP         int            `json:"udp"`          // UDP connections
	Listening   int            `json:"listening"`    // Listening connections
	Established int            `json:"established"`  // Established connections
	ByState     map[string]int `json:"by_state"`     // Connections grouped by state
	ByFamily    map[string]int `json:"by_family"`    // Connections grouped by family (IPv4/IPv6)
}

// GetActiveConnections collects active network connection statistics
func GetActiveConnections() (*ConnectionMetrics, error) {
	connections, err := net.Connections("all") // Get all connection types
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}

	metrics := &ConnectionMetrics{
		ByState:  make(map[string]int),
		ByFamily: make(map[string]int),
	}

	for _, conn := range connections {
		metrics.Total++

		// Count by type
		switch conn.Type {
		case 1: // SOCK_STREAM (TCP)
			metrics.TCP++
		case 2: // SOCK_DGRAM (UDP)
			metrics.UDP++
		}

		// Count by state
		state := strings.ToUpper(conn.Status)
		if state == "" {
			state = "UNKNOWN"
		}
		metrics.ByState[state]++

		// Special state counts
		switch state {
		case "LISTEN":
			metrics.Listening++
		case "ESTABLISHED":
			metrics.Established++
		}

		// Count by family
		switch conn.Family {
		case 2: // AF_INET (IPv4)
			metrics.ByFamily["IPv4"]++
		case 10: // AF_INET6 (IPv6)
			metrics.ByFamily["IPv6"]++
		default:
			metrics.ByFamily["Other"]++
		}
	}

	return metrics, nil
}

// GetListeningPorts returns a list of ports that are currently listening
func GetListeningPorts() ([]int, error) {
	connections, err := net.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get TCP connections: %w", err)
	}

	var listeningPorts []int
	portMap := make(map[uint32]bool) // Use map to avoid duplicates

	for _, conn := range connections {
		if strings.ToUpper(conn.Status) == "LISTEN" && conn.Laddr.Port > 0 {
			if !portMap[conn.Laddr.Port] {
				listeningPorts = append(listeningPorts, int(conn.Laddr.Port))
				portMap[conn.Laddr.Port] = true
			}
		}
	}

	return listeningPorts, nil
}

// GetConnectionsByProcess returns connection counts grouped by process (if available)
func GetConnectionsByProcess() (map[int32]int, error) {
	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}

	processCounts := make(map[int32]int)

	for _, conn := range connections {
		if conn.Pid > 0 {
			processCounts[conn.Pid]++
		}
	}

	return processCounts, nil
}