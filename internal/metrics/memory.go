package metrics

import (
	"log"

	"github.com/DiegoAndradeD/system-monitor/utils"
	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryMetrics struct {
	Usage                 float64
	TotalUsed             float64
	TotalAvailable        float64
	SwapMemoryTotal       float64
	SwapMemoryUsed        float64
	SwapMemoryUsedPercent float64
}

func GetMemoryMetrics() MemoryMetrics {
	metrics := MemoryMetrics{}

	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting virtual memory stats: %v", err)
		return metrics
	}

	metrics.Usage = virtualMemory.UsedPercent
	metrics.TotalUsed = utils.ConvertBytesToGb(virtualMemory.Used)
	metrics.TotalAvailable = utils.ConvertBytesToGb(virtualMemory.Available)

	swapMemory, err := mem.SwapMemory()
	if err != nil {
		log.Printf("Error getting swap memory stats: %v", err)
		return metrics
	}

	metrics.SwapMemoryTotal = utils.ConvertBytesToGb(swapMemory.Total)
	metrics.SwapMemoryUsed = utils.ConvertBytesToGb(swapMemory.Used)
	metrics.SwapMemoryUsedPercent = swapMemory.UsedPercent

	return metrics
}
