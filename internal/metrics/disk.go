package metrics

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
)

// DiskMetrics represents disk usage for a single partition
type DiskMetrics struct {
	Partition   string  `json:"partition"`    // Mount point or drive letter
	Filesystem  string  `json:"filesystem"`   // Filesystem type
	Total       uint64  `json:"total"`        // Total space in bytes
	Used        uint64  `json:"used"`         // Used space in bytes
	Free        uint64  `json:"free"`         // Free space in bytes
	UsedPercent float64 `json:"used_percent"` // Used percentage
}

// GetDiskUsage collects disk usage metrics for all mounted partitions
func GetDiskUsage() ([]*DiskMetrics, error) {
	partitions, err := disk.Partitions(false) // false = exclude virtual filesystems
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %w", err)
	}

	var diskMetrics []*DiskMetrics

	for _, partition := range partitions {
		// Skip certain filesystem types that we don't want to monitor
		if shouldSkipFilesystem(partition.Fstype) {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			// Log error but continue with other partitions
			continue
		}

		diskMetrics = append(diskMetrics, &DiskMetrics{
			Partition:   partition.Mountpoint,
			Filesystem:  partition.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	if len(diskMetrics) == 0 {
		return nil, fmt.Errorf("no disk partitions found")
	}

	return diskMetrics, nil
}

// GetDiskIO collects disk I/O statistics (optional)
func GetDiskIO() (map[string]interface{}, error) {
	ioStats, err := disk.IOCounters()
	if err != nil {
		// I/O stats might not be available on all systems
		return map[string]interface{}{
			"read_bytes":  uint64(0),
			"write_bytes": uint64(0),
			"read_count":  uint64(0),
			"write_count": uint64(0),
		}, nil
	}

	var totalReadBytes, totalWriteBytes, totalReadCount, totalWriteCount uint64

	for _, stat := range ioStats {
		totalReadBytes += stat.ReadBytes
		totalWriteBytes += stat.WriteBytes
		totalReadCount += stat.ReadCount
		totalWriteCount += stat.WriteCount
	}

	return map[string]interface{}{
		"read_bytes":  totalReadBytes,
		"write_bytes": totalWriteBytes,
		"read_count":  totalReadCount,
		"write_count": totalWriteCount,
	}, nil
}

// shouldSkipFilesystem determines if a filesystem should be skipped
func shouldSkipFilesystem(fstype string) bool {
	skipTypes := []string{
		"tmpfs", "devtmpfs", "sysfs", "proc", "devpts",
		"cgroup", "cgroup2", "pstore", "bpf", "tracefs",
		"debugfs", "mqueue", "hugetlbfs", "systemd-1",
		"binfmt_misc", "autofs", "rpc_pipefs", "nfsd",
	}

	fstype = strings.ToLower(fstype)
	for _, skipType := range skipTypes {
		if fstype == skipType {
			return true
		}
	}

	return false
}