package metrics

import (
	"log"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryMetrics struct {
	Usage float64
}

func GetMemoryMetrics() MemoryMetrics {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting virtual memory stats: %v", err)
		return MemoryMetrics{Usage: 0}
	}
	return MemoryMetrics{Usage: virtualMemory.UsedPercent}
}
