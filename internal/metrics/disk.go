package metrics

import (
	"log"

	"github.com/shirou/gopsutil/disk"
	"golang.org/x/exp/slices"
)

type DiskMetrics struct {
	Usage          float64
	TotalUsed      float64
	TotalAvailable float64
	TotalSize      float64
}

var nonPhysicalPartitionsTypes = []string{"tmpfs", "overlay", "proc", "sysfs", "cgroup", "squashfs", "devtmpfs", "vfat"}

func GetDiskMetrics() DiskMetrics {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Printf("Error getting disk partitions: %v", err)
		return DiskMetrics{}
	}

	var totalUsed float64 = 0
	var totalSize float64 = 0

	for _, p := range partitions {
		if slices.Contains(nonPhysicalPartitionsTypes, p.Fstype) {
			continue
		}

		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			log.Printf("Error getting usage for partition %v: %v", p.Mountpoint, err)
			continue
		}

		totalUsed += float64(usage.Used)
		totalSize += float64(usage.Total)
	}

	if totalSize == 0 {
		return DiskMetrics{}
	}

	totalAvailable := totalSize - totalUsed

	return DiskMetrics{
		Usage:          (totalUsed / totalSize) * 100,
		TotalUsed:      totalUsed,
		TotalAvailable: totalAvailable,
		TotalSize:      totalSize,
	}
}
