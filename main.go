package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemMetrics struct {
	CPUUsage    float64
	MemoryUsage float64
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

func (m *Monitor) Start() {
	go m.collectorRoutine(getCPUPercentage, m.updateCPUUsage)
	go m.collectorRoutine(getMemoryPercentage, m.updateMemoryUsage)
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

func getCPUPercentage() float64 {
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

func getMemoryPercentage() float64 {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting virtual memory stats: %v", err)
		return 0
	}
	return virtualMemory.UsedPercent
}

func main() {
	monitor := NewMonitor()
	monitor.Start()
	defer monitor.Stop()

	const (
		screenWidth  = 800
		screenHeight = 600
	)

	rl.InitWindow(screenWidth, screenHeight, "Go System Monitor - Refatorado")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		metrics := monitor.GetMetrics()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.DrawText("System Monitor", 10, 10, 40, rl.LightGray)
		rl.DrawText(fmt.Sprintf("CPU USAGE: %.2f%%", metrics.CPUUsage), 10, 100, 20, rl.White)
		rl.DrawText(fmt.Sprintf("MEMORY USAGE: %.2f%%", metrics.MemoryUsage), 10, 150, 20, rl.White)

		rl.EndDrawing()
	}
}
