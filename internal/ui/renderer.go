package ui

import (
	"fmt"

	"github.com/DiegoAndradeD/system-monitor/internal/metrics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MetricDisplay struct {
	Label string
	Value float64
	Unit  string
}

type MetricsSection struct {
	Title   string
	Color   rl.Color
	Metrics []MetricDisplay
}

func Render(metrics metrics.SystemMetrics) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	rl.DrawText("System Monitor", 10, 10, 40, rl.LightGray)

	yOffset := 80

	cpuSection := MetricsSection{
		Title: "CPU",
		Color: rl.SkyBlue,
		Metrics: []MetricDisplay{
			{Label: "Usage", Value: metrics.CPUMetrics.Usage, Unit: "%"},
			{Label: "Frequency", Value: metrics.CPUMetrics.Frequency, Unit: "MHz"},
		},
	}

	DrawMetricsSection(cpuSection, float32(yOffset))

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

func DrawMetricsSection(section MetricsSection, yPosition float32) {
	screenW := float32(rl.GetScreenWidth())

	const (
		rectWidthRatio = 0.9
		padding        = 20
		fontSize       = 20
		lineHeight     = 60
	)

	rectWidth := screenW * rectWidthRatio
	rectX := (screenW - rectWidth) / 2
	rectHeight := float32(len(section.Metrics)*int(lineHeight) + 80)

	rl.DrawRectangleLines(int32(rectX), int32(yPosition), int32(rectWidth), int32(rectHeight), rl.Gray)
	rl.DrawText(section.Title, int32(rectX)+padding, int32(yPosition)+padding, fontSize+5, section.Color)
	currentY := yPosition + padding + 40
	for _, metric := range section.Metrics {
		DrawPercentageBars(metric.Value, rectX+padding+300, currentY)

		labelText := metric.Label + ":"
		rl.DrawText(labelText, int32(rectX)+padding, int32(currentY), fontSize, rl.White)

		valueText := fmt.Sprintf("%.2f %s", metric.Value, metric.Unit)
		rl.DrawText(valueText, int32(rectX)+padding+150, int32(currentY), fontSize, rl.LightGray)

		currentY += lineHeight
	}
}
