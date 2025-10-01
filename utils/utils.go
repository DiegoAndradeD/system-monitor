package utils

import "fmt"

type SingleValue float64

type PercentageValue float64

type DualValue struct {
	Used      float64
	Available float64
}

func ConvertBytesToGb(value uint64) float64 {
	return float64(value / 1024 / 1024 / 1024)
}

func (v SingleValue) Format(unit string) string {
	return fmt.Sprintf("%.2f %s", v, unit)
}

func (v DualValue) Format(unit string) string {
	return fmt.Sprintf("%.2f / %.2f %s", v.Used, v.Available, unit)
}

func (v PercentageValue) Format(unit string) string {
	return fmt.Sprintf("%.2f%s", v, unit)
}
