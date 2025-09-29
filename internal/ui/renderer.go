package ui

import (
	"github.com/DiegoAndradeD/system-monitor/internal/metrics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func Render(metrics metrics.SystemMetrics) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	rl.DrawText("System Monitor", 10, 10, 40, rl.LightGray)
	DrawMetricsSection(metrics.CPUUsage, "CPU:", rl.White)

	rl.EndDrawing()
}

func DrawPercentageBars(percent float64, x, y float32) {
	const (
		totalBars = 50
		barWidth  = 10
		barHeight = 20
		gap       = 2
	)

	filledBars := int(percent / 100 * totalBars)

	for i := range totalBars {
		color := rl.DarkGray
		if i < filledBars {
			color = rl.SkyBlue
		}
		rl.DrawRectangle(
			int32(x)+int32(i*(barWidth+gap)),
			int32(y),
			int32(barWidth),
			int32(barHeight),
			color,
		)
	}
}

func DrawMetricsSection(metric float64, title string, color rl.Color) {
	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())

	const (
		rectWidthRatio = 0.8
		rectHeight     = 150
		padding        = 20
		fontSize       = 20
	)

	rectWidth := screenW * rectWidthRatio
	rectX := (screenW - rectWidth) / 2
	rectY := (screenH - rectHeight) / 2

	rl.DrawRectangleLines(int32(rectX), int32(rectY), int32(rectWidth), int32(rectHeight), rl.Gray)
	rl.DrawText(title, int32(rectX)+padding, int32(rectY)+padding, fontSize, color)
	DrawPercentageBars(metric, rectX+padding, rectY+padding*3)
}
