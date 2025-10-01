package ui

import (
	"github.com/DiegoAndradeD/system-monitor/internal/metrics"
	"github.com/DiegoAndradeD/system-monitor/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenPadding = 10
	TitleFontSize = 40
	TitleOffsetY  = 10
	MetricsStartY = 80

	SectionRectWidthRatio = 0.9
	SectionPadding        = 20
	SectionTitleFontSize  = 25
	SectionTitleOffsetY   = 20
	SectionMetricsStartY  = 60

	MetricFontSize   = 20
	MetricLineHeight = 30
	MetricLabelWidth = 125
	MetricValueWidth = 400

	BarCount   = 50
	BarWidth   = 8
	BarHeight  = 16
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
	Title   string
	Color   rl.Color
	Metrics []MetricDisplay
}

func Render(sections []MetricsSection) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.Black)

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

func CreateNetworkSection(netMetrics metrics.NetworkMetrics) MetricsSection {
	return MetricsSection{
		Title: "Network",
		Color: rl.Violet,
		Metrics: []MetricDisplay{
			{Label: "Upload", Value: utils.SingleValue(netMetrics.UploadSpeedKBps), Unit: "KB/s"},
			{Label: "Download", Value: utils.SingleValue(netMetrics.DownloadSpeedKBps), Unit: "KB/s"},
			{Label: "Total Sent", Value: utils.SingleValue(netMetrics.TotalSentGB), Unit: "GB"},
			{Label: "Total Recv", Value: utils.SingleValue(netMetrics.TotalRecvGB), Unit: "GB"},
		},
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

	sectionHeight := SectionMetricsStartY + float32(len(section.Metrics))*MetricLineHeight + SectionPadding

	rl.DrawRectangleLines(
		int32(rectX),
		int32(yPosition),
		int32(rectWidth),
		int32(sectionHeight),
		rl.Gray,
	)

	rl.DrawText(
		section.Title,
		int32(rectX)+SectionPadding,
		int32(yPosition)+SectionTitleOffsetY,
		SectionTitleFontSize,
		rl.White,
	)

	currentY := yPosition + SectionMetricsStartY
	for _, metric := range section.Metrics {
		drawMetric(metric, section.Color, rectX, currentY)
		currentY += MetricLineHeight
	}

	return sectionHeight
}

func drawMetric(metric MetricDisplay, sectionColor rl.Color, rectX, y float32) {
	baseX := int32(rectX) + SectionPadding
	var valuePadding int32 = 0

	if percentageValue, ok := metric.Value.(utils.PercentageValue); ok {
		drawPercentageBars(float64(percentageValue), rectX+SectionPadding, y, sectionColor)
		valuePadding = 100
	}

	labelX := baseX
	valueX := baseX + MetricLabelWidth

	if _, ok := metric.Value.(utils.PercentageValue); ok {
		labelX += BarOffsetX
		valueX += BarOffsetX
	}

	if metric.Label != "" {
		rl.DrawText(metric.Label+":", labelX, int32(y), MetricFontSize, rl.White)
	}

	valueText := metric.Value.Format(metric.Unit)
	rl.DrawText(valueText, valueX+valuePadding, int32(y), MetricFontSize, rl.LightGray)
}

func RenderSystemMetrics(m metrics.SystemMetrics) {
	sections := []MetricsSection{
		CreateCPUSection(m.CPUMetrics),
		CreateMemorySection(m.MemoryMetrics),
		CreateDiskSection(m.DiskMetrics),
		CreateNetworkSection(m.NetworkMetrics),
	}
	Render(sections)
}
