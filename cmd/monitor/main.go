package main

import (
	"github.com/DiegoAndradeD/system-monitor/internal/metrics"
	"github.com/DiegoAndradeD/system-monitor/internal/ui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 800
	TargetFPS    = 60
	WindowTitle  = "Go System Monitor"
)

func main() {
	monitor := metrics.NewMonitor()
	monitor.Start()
	defer monitor.Stop()

	rl.InitWindow(ScreenWidth, ScreenHeight, WindowTitle)
	defer rl.CloseWindow()
	rl.SetTargetFPS(TargetFPS)

	for !rl.WindowShouldClose() {
		ui.RenderSystemMetrics(monitor.GetMetrics())
	}
}
