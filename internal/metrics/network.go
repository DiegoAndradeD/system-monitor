package metrics

import (
	"log"
	"time"

	"github.com/DiegoAndradeD/system-monitor/utils"
	"github.com/shirou/gopsutil/v3/net"
)

type NetworkMetrics struct {
	UploadSpeedKBps   float64
	DownloadSpeedKBps float64
	TotalSentGB       float64
	TotalRecvGB       float64
}

type NetworkCollector struct {
	prevStats net.IOCountersStat
	prevTime  time.Time
}

func NewNetworkCollector() *NetworkCollector {
	stats, err := getSystemWideIOCounters()
	if err != nil {
		return &NetworkCollector{
			prevStats: net.IOCountersStat{},
			prevTime:  time.Now(),
		}
	}
	return &NetworkCollector{
		prevStats: stats,
		prevTime:  time.Now(),
	}
}

func (nc *NetworkCollector) GetNetworkMetrics() NetworkMetrics {
	currentStats, err := getSystemWideIOCounters()
	if err != nil {
		log.Printf("Error getting network I/O counters: %v", err)
		return NetworkMetrics{}
	}
	currentTime := time.Now()

	duration := currentTime.Sub(nc.prevTime).Seconds()
	if duration <= 0 {
		return NetworkMetrics{}
	}

	bytesSent := currentStats.BytesSent - nc.prevStats.BytesSent
	bytesRecv := currentStats.BytesRecv - nc.prevStats.BytesRecv

	uploadSpeed := float64(bytesSent) / duration / 1024
	downloadSpeed := float64(bytesRecv) / duration / 1024

	nc.prevStats = currentStats
	nc.prevTime = currentTime

	return NetworkMetrics{
		UploadSpeedKBps:   uploadSpeed,
		DownloadSpeedKBps: downloadSpeed,
		TotalSentGB:       utils.ConvertBytesToGb(currentStats.BytesSent),
		TotalRecvGB:       utils.ConvertBytesToGb(currentStats.BytesRecv),
	}
}

func getSystemWideIOCounters() (net.IOCountersStat, error) {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return net.IOCountersStat{}, err
	}
	return stats[0], nil
}
