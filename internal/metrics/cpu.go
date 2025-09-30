package metrics

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUMetrics struct {
	Usage     float64
	Frequency float64
}

func GetCPUMetrics() CPUMetrics {
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Printf("Error getting CPU info: %v", err)
		return CPUMetrics{Usage: 0, Frequency: 0}
	}
	maxFrequency := cpuInfo[0].Mhz

	cpuPercents, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Printf("Error getting CPU Usage: %v", err)
		return CPUMetrics{Usage: 0, Frequency: 0}
	}

	var dynamicFrequency float64

	for _, percent := range cpuPercents {
		dynamicFrequency = maxFrequency * (percent / 100.0)
	}

	if len(cpuPercents) > 0 {
		return CPUMetrics{Usage: cpuPercents[0], Frequency: dynamicFrequency}
	}
	return CPUMetrics{Usage: 0, Frequency: 0}

}
