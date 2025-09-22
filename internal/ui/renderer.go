package ui

import (
	"fmt"

	"github.com/DiegoAndradeD/system-monitor/internal/metrics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func Render(metrics metrics.SystemMetrics) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	rl.DrawText("System Monitor", 10, 10, 40, rl.LightGray)
	rl.DrawText(fmt.Sprintf("CPU USAGE: %.2f%%", metrics.CPUUsage), 10, 100, 20, rl.White)
	rl.DrawText(fmt.Sprintf("MEMORY USAGE: %.2f%%", metrics.MemoryUsage), 10, 150, 20, rl.White)
	rl.DrawText(fmt.Sprintf("DISK USAGE: %.2F%%", metrics.MainDiskUsage), 10, 200, 20, rl.White)

	rl.EndDrawing()
}
