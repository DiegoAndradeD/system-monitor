package metrics

import (
	"context"
	"sync"
	"time"
)

type SystemMetrics struct {
	CPUMetrics     CPUMetrics
	MemoryMetrics  MemoryMetrics
	DiskMetrics    DiskMetrics
	NetworkMetrics NetworkMetrics
}

type Monitor struct {
	metrics          SystemMetrics
	networkCollector *NetworkCollector
	mu               sync.RWMutex
	ctx              context.Context
	cancel           context.CancelFunc
}

func NewMonitor() *Monitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Monitor{
		ctx:              ctx,
		cancel:           cancel,
		networkCollector: NewNetworkCollector(),
	}
}

func (m *Monitor) GetMetrics() SystemMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.metrics
}

func (m *Monitor) Start() {
	go collectorRoutine(m, GetCPUMetrics, m.updateCPUMetrics)
	go collectorRoutine(m, GetMemoryMetrics, m.updateMemoryUsage)
	go collectorRoutine(m, GetDiskMetrics, m.updateMainDiskMetrics)
	go collectorRoutine(m, m.networkCollector.GetNetworkMetrics, m.updateNetworkMetrics)

}
func (m *Monitor) Stop() {
	m.cancel()
}

func collectorRoutine[T any](m *Monitor, getData func() T, updateData func(T)) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	updateData(getData())

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			updateData(getData())
		}
	}
}

func (m *Monitor) updateCPUMetrics(metrics CPUMetrics) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics.CPUMetrics = CPUMetrics{Usage: metrics.Usage, Frequency: metrics.Frequency}
}

func (m *Monitor) updateMemoryUsage(metrics MemoryMetrics) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics.MemoryMetrics = MemoryMetrics{Usage: metrics.Usage, TotalUsed: metrics.TotalUsed,
		TotalAvailable: metrics.TotalAvailable, SwapMemoryTotal: metrics.SwapMemoryTotal,
		SwapMemoryUsed: metrics.SwapMemoryUsed, SwapMemoryUsedPercent: metrics.SwapMemoryUsedPercent}
}

func (m *Monitor) updateMainDiskMetrics(metrics DiskMetrics) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics.DiskMetrics = DiskMetrics{
		Usage:          metrics.Usage,
		TotalUsed:      metrics.TotalUsed,
		TotalAvailable: metrics.TotalAvailable,
		TotalSize:      metrics.TotalSize,
	}
}

func (m *Monitor) updateNetworkMetrics(metrics NetworkMetrics) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics.NetworkMetrics = metrics
}
