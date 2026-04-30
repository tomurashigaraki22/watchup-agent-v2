package metrics

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/net"
)

// NetworkMetrics represents network I/O for a single interface
type NetworkMetrics struct {
	Interface   string `json:"interface"`    // Interface name
	BytesSent   uint64 `json:"bytes_sent"`   // Bytes sent
	BytesRecv   uint64 `json:"bytes_recv"`   // Bytes received
	PacketsSent uint64 `json:"packets_sent"` // Packets sent
	PacketsRecv uint64 `json:"packets_recv"` // Packets received
	Errin       uint64 `json:"errin"`        // Input errors
	Errout      uint64 `json:"errout"`       // Output errors
	Dropin      uint64 `json:"dropin"`       // Input drops
	Dropout     uint64 `json:"dropout"`      // Output drops
}

// NetworkSummary represents aggregated network metrics
type NetworkSummary struct {
	TotalBytesSent   uint64 `json:"total_bytes_sent"`
	TotalBytesRecv   uint64 `json:"total_bytes_recv"`
	TotalPacketsSent uint64 `json:"total_packets_sent"`
	TotalPacketsRecv uint64 `json:"total_packets_recv"`
	ActiveInterfaces int    `json:"active_interfaces"`
}

// GetNetworkUsage collects network I/O metrics for all interfaces
func GetNetworkUsage() ([]*NetworkMetrics, error) {
	ioStats, err := net.IOCounters(true) // true = per interface
	if err != nil {
		return nil, fmt.Errorf("failed to get network I/O stats: %w", err)
	}

	var networkMetrics []*NetworkMetrics

	for _, stat := range ioStats {
		// Skip loopback and inactive interfaces
		if shouldSkipInterface(stat.Name) {
			continue
		}

		networkMetrics = append(networkMetrics, &NetworkMetrics{
			Interface:   stat.Name,
			BytesSent:   stat.BytesSent,
			BytesRecv:   stat.BytesRecv,
			PacketsSent: stat.PacketsSent,
			PacketsRecv: stat.PacketsRecv,
			Errin:       stat.Errin,
			Errout:      stat.Errout,
			Dropin:      stat.Dropin,
			Dropout:     stat.Dropout,
		})
	}

	return networkMetrics, nil
}

// GetNetworkSummary returns aggregated network metrics
func GetNetworkSummary() (*NetworkSummary, error) {
	interfaces, err := GetNetworkUsage()
	if err != nil {
		return nil, err
	}

	summary := &NetworkSummary{}

	for _, iface := range interfaces {
		summary.TotalBytesSent += iface.BytesSent
		summary.TotalBytesRecv += iface.BytesRecv
		summary.TotalPacketsSent += iface.PacketsSent
		summary.TotalPacketsRecv += iface.PacketsRecv
		summary.ActiveInterfaces++
	}

	return summary, nil
}

// shouldSkipInterface determines if a network interface should be skipped
func shouldSkipInterface(name string) bool {
	name = strings.ToLower(name)
	
	// Skip loopback interfaces
	if strings.Contains(name, "lo") || strings.Contains(name, "loopback") {
		return true
	}

	// Skip virtual interfaces (common patterns)
	skipPatterns := []string{
		"docker", "veth", "br-", "virbr", "vmnet", "vbox",
		"tun", "tap", "wg", "ppp", "slip",
	}

	for _, pattern := range skipPatterns {
		if strings.Contains(name, pattern) {
			return true
		}
	}

	return false
}