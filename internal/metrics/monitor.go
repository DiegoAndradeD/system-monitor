package metrics

import (
	"context"
	"sync"
	"time"
)

type SystemMetrics struct {
	CPUMetrics    CPUMetrics
	MemoryUsage   MemoryMetrics
	MainDiskUsage DiskMetrics
}

type Monitor struct {
	metrics SystemMetrics
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewMonitor() *Monitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Monitor{
		ctx:    ctx,
		cancel: cancel,
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
	go collectorRoutine(m, GetDiskMetrics, m.updateMainDiskUsage)
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
	m.metrics.MemoryUsage = MemoryMetrics{Usage: metrics.Usage}
}

func (m *Monitor) updateMainDiskUsage(metrics DiskMetrics) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics.MainDiskUsage = DiskMetrics{Usage: metrics.Usage}
}
