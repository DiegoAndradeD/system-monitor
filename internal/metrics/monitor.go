package metrics

import (
	"context"
	"sync"
	"time"
)

type SystemMetrics struct {
	CPUUsage      float64
	MemoryUsage   float64
	MainDiskUsage float64
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
	go m.collectorRoutine(GetCPUPercentage, m.updateCPUUsage)
	go m.collectorRoutine(GetMemoryPercentage, m.updateMemoryUsage)
	go m.collectorRoutine(GetDiskUsage, m.updateMainDiskUsage)
}

func (m *Monitor) Stop() {
	m.cancel()
}

func (m *Monitor) collectorRoutine(getData func() float64, updateData func(float64)) {
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

func (m *Monitor) updateCPUUsage(usage float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics.CPUUsage = usage
}

func (m *Monitor) updateMemoryUsage(usage float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics.MemoryUsage = usage
}

func (m *Monitor) updateMainDiskUsage(usage float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics.MainDiskUsage = usage
}
