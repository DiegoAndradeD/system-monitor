package main

import (
	"github.com/DiegoAndradeD/system-monitor/internal/metrics"
	"github.com/DiegoAndradeD/system-monitor/internal/ui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

func main() {

	monitor := metrics.NewMonitor()
	monitor.Start()
	defer monitor.Stop()

	rl.InitWindow(ScreenWidth, ScreenHeight, "Go System Monitor")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		ui.Render(monitor.GetMetrics())
	}
}
