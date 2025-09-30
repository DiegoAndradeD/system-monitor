package metrics

import (
	"log"

	"github.com/shirou/gopsutil/disk"
	"golang.org/x/exp/slices"
)

type DiskMetrics struct {
	Usage float64
}

var nonPhysicalPartitionsTypes = []string{"tmpfs", "overlay", "proc", "sysfs", "cgroup", "squashfs", "devtmpfs", "vfat"}

func GetDiskMetrics() DiskMetrics {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Printf("Error getting disk partitions: %v", err)
		return DiskMetrics{Usage: 0}
	}

	var totalUsed uint64 = 0
	var totalSize uint64 = 0

	for _, p := range partitions {
		if slices.Contains(nonPhysicalPartitionsTypes, p.Fstype) {
			continue
		}

		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			log.Printf("Error getting usage for partition %v: %v", p.Mountpoint, err)
			continue
		}

		totalUsed += usage.Used
		totalSize += usage.Total
	}

	if totalSize == 0 {
		return DiskMetrics{Usage: 0}
	}

	return DiskMetrics{Usage: (float64(totalUsed) / float64(totalSize)) * 100}
}
