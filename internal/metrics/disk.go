package metrics

import (
	"log"

	"github.com/shirou/gopsutil/disk"
	"golang.org/x/exp/slices"
)

var nonPhysicalPartitionsTypes = []string{"tmpfs", "overlay", "proc", "sysfs", "cgroup", "squashfs", "devtmpfs", "vfat"}

func GetDiskUsage() float64 {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Printf("Error getting disk partitions: %v", err)
		return 0
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
		return 0
	}

	return (float64(totalUsed) / float64(totalSize)) * 100
}
