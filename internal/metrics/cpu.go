package metrics

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

func GetCPUPercentage() float64 {
	cpuUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
		return 0
	}
	if len(cpuUsage) > 0 {
		return cpuUsage[0]
	}
	return 0
}
