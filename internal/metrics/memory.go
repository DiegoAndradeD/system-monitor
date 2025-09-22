package metrics

import (
	"log"

	"github.com/shirou/gopsutil/v3/mem"
)

func GetMemoryPercentage() float64 {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting virtual memory stats: %v", err)
		return 0
	}
	return virtualMemory.UsedPercent
}
