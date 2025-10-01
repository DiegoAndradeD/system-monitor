package ui

import (
	"fmt"

	"github.com/DiegoAndradeD/system-monitor/internal/metrics"
	"github.com/DiegoAndradeD/system-monitor/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenPadding = 10
	TitleFontSize = 40
	TitleOffsetY  = 10
	MetricsStartY = 50

	SectionRectWidthRatio = 0.9
	SectionPadding        = 15
	SectionTitleFontSize  = 22
	SectionTitleOffsetY   = 10
	SectionMetricsStartY  = 45

	MetricFontSize   = 18
	MetricLineHeight = 25
	MetricLabelWidth = 125
	MetricValueWidth = 400

	BarCount   = 50
	BarWidth   = 8
	BarHeight  = 14
	BarGap     = 2
	BarOffsetX = 300
)

type MetricValue interface {
	Format(unit string) string
}

type MetricDisplay struct {
	Label string
	Value MetricValue
	Unit  string
}

type MetricsSection struct {
	Title           string
	Color           rl.Color
	Metrics         []MetricDisplay
	UploadHistory   []float64
	DownloadHistory []float64
}

func Render(sections []MetricsSection) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.Black)

	rl.DrawRectangleGradientV(
		0, 0,
		int32(rl.GetScreenWidth()), 700,
		rl.NewColor(20, 20, 40, 100),
		rl.NewColor(10, 10, 20, 120),
	)

	yOffset := MetricsStartY
	for _, section := range sections {
		height := drawMetricsSection(section, float32(yOffset))
		yOffset += int(height) + SectionPadding
	}
}

func CreateCPUSection(cpuMetrics metrics.CPUMetrics) MetricsSection {
	return MetricsSection{
		Title: "CPU",
		Color: rl.SkyBlue,
		Metrics: []MetricDisplay{
			{Label: "", Value: utils.PercentageValue(cpuMetrics.Usage), Unit: "%"},
			{Label: "Frequency", Value: utils.SingleValue(cpuMetrics.Frequency), Unit: "MHz"},
		},
	}
}

func CreateMemorySection(memMetrics metrics.MemoryMetrics) MetricsSection {
	return MetricsSection{
		Title: "Memory",
		Color: rl.Green,
		Metrics: []MetricDisplay{
			{Label: "", Value: utils.PercentageValue(memMetrics.Usage), Unit: "%"},
			{Label: "Used", Value: utils.DualValue{Used: memMetrics.TotalUsed, Available: memMetrics.TotalAvailable}, Unit: "GB"},
			{Label: "Swap", Value: utils.DualValue{Used: memMetrics.SwapMemoryUsed, Available: memMetrics.SwapMemoryTotal}, Unit: "GB"},
		},
	}
}

func CreateDiskSection(diskMetrics metrics.DiskMetrics) MetricsSection {
	return MetricsSection{
		Title: "Disk",
		Color: rl.Orange,
		Metrics: []MetricDisplay{
			{
				Label: "",
				Value: utils.PercentageValue(diskMetrics.Usage),
				Unit:  "%",
			},
			{
				Label: "Used",
				Value: utils.DualValue{
					Used:      diskMetrics.TotalUsed / (1024 * 1024 * 1024),
					Available: diskMetrics.TotalSize / (1024 * 1024 * 1024),
				},
				Unit: "GB",
			},
		},
	}
}

func CreateNetworkSection(netMetrics metrics.SystemMetrics) MetricsSection {
	return MetricsSection{
		Title: "Network",
		Color: rl.Violet,
		Metrics: []MetricDisplay{
			{Label: "Upload", Value: utils.SingleValue(netMetrics.NetworkMetrics.UploadSpeedKBps), Unit: "KB/s"},
			{Label: "Download", Value: utils.SingleValue(netMetrics.NetworkMetrics.DownloadSpeedKBps), Unit: "KB/s"},
		},
		UploadHistory:   netMetrics.UploadHistory,
		DownloadHistory: netMetrics.DownloadHistory,
	}
}

func drawPercentageBars(percent float64, x, y float32, color rl.Color) {
	filledBars := int(percent / 100.0 * BarCount)

	for i := range BarCount {
		barColor := rl.DarkGray
		if i < filledBars {
			barColor = color
		}

		barX := int32(x) + int32(i*(BarWidth+BarGap))
		rl.DrawRectangle(barX, int32(y), BarWidth, BarHeight, barColor)
	}
}

func drawMetricsSection(section MetricsSection, yPosition float32) float32 {
	screenW := float32(rl.GetScreenWidth())
	rectWidth := screenW * SectionRectWidthRatio
	rectX := (screenW - rectWidth) / 2

	var graphHeight float32 = 40
	var graphPadding float32 = 5
	metricsTextYOffset := SectionMetricsStartY

	totalHeight := float32(len(section.Metrics))*MetricLineHeight + float32(metricsTextYOffset) + float32(SectionPadding)

	if section.Title == "Network" {
		totalHeight += (graphHeight * 2) + graphPadding
	}

	rl.DrawRectangleLines(
		int32(rectX),
		int32(yPosition),
		int32(rectWidth),
		int32(totalHeight),
		rl.Gray,
	)

	rl.DrawText(
		section.Title,
		int32(rectX)+SectionPadding,
		int32(yPosition)+SectionTitleOffsetY,
		SectionTitleFontSize,
		rl.White,
	)

	currentY := yPosition

	if section.Title == "Network" {
		graphY := yPosition + float32(metricsTextYOffset)
		graphWidth := rectWidth - float32(SectionPadding*2)
		drawMetricLineGraph(section.UploadHistory, rectX+float32(SectionPadding), graphY, graphWidth, graphHeight, section.Color)
		graphY += graphHeight + graphPadding
		drawMetricLineGraph(section.DownloadHistory, rectX+float32(SectionPadding), graphY, graphWidth, graphHeight, rl.SkyBlue)
		currentY += (graphHeight * 2.2) + graphPadding
	}

	currentY += float32(metricsTextYOffset)

	for _, metric := range section.Metrics {
		drawMetric(metric, section.Color, rectX, currentY)
		currentY += MetricLineHeight
	}

	return totalHeight
}

func drawMetric(metric MetricDisplay, sectionColor rl.Color, rectX, y float32) {
	baseX := int32(rectX) + SectionPadding
	labelX := baseX
	valueX := baseX + MetricLabelWidth
	var valuePadding int32 = 0

	if percentageValue, ok := metric.Value.(utils.PercentageValue); ok {
		drawPercentageBars(float64(percentageValue), float32(baseX), y, sectionColor)
		labelX += BarOffsetX
		valueX += BarOffsetX
		valuePadding = 100
	}

	if metric.Label != "" {
		rl.DrawText(metric.Label+":", labelX, int32(y), MetricFontSize, rl.White)
	}

	valueText := metric.Value.Format(metric.Unit)
	rl.DrawText(valueText, valueX+valuePadding, int32(y), MetricFontSize, rl.LightGray)
}

func drawMetricLineGraph(history []float64, x, y, width, height float32, color rl.Color) {
	if len(history) < 2 {
		return
	}

	maxValue := 0.0
	for _, v := range history {
		if v > maxValue {
			maxValue = v
		}
	}
	if maxValue < 10 {
		maxValue = 10
	}

	rl.DrawRectangleLines(int32(x), int32(y), int32(width), int32(height), rl.DarkGray)

	stepX := width / float32(len(history)-1)
	for i := 1; i < len(history); i++ {
		x1 := x + float32(i-1)*stepX
		y1 := y + height - (float32(history[i-1])/float32(maxValue))*height

		x2 := x + float32(i)*stepX
		y2 := y + height - (float32(history[i])/float32(maxValue))*height

		if y1 < y {
			y1 = y
		}
		if y2 < y {
			y2 = y
		}

		rl.DrawLineEx(rl.NewVector2(x1, y1), rl.NewVector2(x2, y2), 2.0, color)
	}

	currentValueText := utils.SingleValue(history[len(history)-1]).Format("KB/s")
	peakValueText := fmt.Sprintf("Peak: %.1f KB/s", maxValue)
	rl.DrawText(currentValueText, int32(x)+5, int32(y)+5, 12, rl.White)
	rl.DrawText(
		peakValueText,
		int32(x+width-float32(rl.MeasureText(peakValueText, 12)))-5,
		int32(y)+5,
		12,
		rl.LightGray,
	)
}

func RenderSystemMetrics(m metrics.SystemMetrics) {
	sections := []MetricsSection{
		CreateCPUSection(m.CPUMetrics),
		CreateMemorySection(m.MemoryMetrics),
		CreateDiskSection(m.DiskMetrics),
		CreateNetworkSection(m),
	}
	Render(sections)
}
