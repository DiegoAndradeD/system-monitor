package main

import (
	"github.com/DiegoAndradeD/system-monitor/internal/metrics"
	"github.com/DiegoAndradeD/system-monitor/internal/ui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	monitor := metrics.NewMonitor()
	monitor.Start()
	defer monitor.Stop()

	const (
		screenWidth  = 800
		screenHeight = 600
	)

	rl.InitWindow(screenWidth, screenHeight, "Go System Monitor")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		ui.Render(monitor.GetMetrics())
	}
}
