package ui

import (
	"fmt"

	"fyne.io/fyne/v2/canvas"
)

// MetricLabels содержит все метки для отображения метрик
type MetricLabels struct {
	CPUTemp  *canvas.Text
	CPUUsage *canvas.Text
	GPUTemp  *canvas.Text
	GPUUsage *canvas.Text
	RAMUsage *canvas.Text
	HDDTemp  *canvas.Text
}

// NewMetricLabels создает новый набор меток для метрик
func NewMetricLabels() *MetricLabels {
	return &MetricLabels{
		CPUTemp:  NewValueLabel(32),
		CPUUsage: NewValueLabel(32),
		GPUTemp:  NewValueLabel(32),
		GPUUsage: NewValueLabel(32),
		RAMUsage: NewValueLabel(32),
		HDDTemp:  NewValueLabel(32),
	}
}

// UpdateMetrics обновляет значения метрик
func (m *MetricLabels) UpdateMetrics(cpuTemp, cpuUsage, gpuTemp, gpuUsage, ramUsage, hddTemp float64) {
	m.CPUTemp.Text = formatMetricValue(cpuTemp)
	canvas.Refresh(m.CPUTemp)

	m.CPUUsage.Text = formatMetricValue(cpuUsage)
	canvas.Refresh(m.CPUUsage)

	m.GPUTemp.Text = formatMetricValue(gpuTemp)
	canvas.Refresh(m.GPUTemp)

	m.GPUUsage.Text = formatMetricValue(gpuUsage)
	canvas.Refresh(m.GPUUsage)

	m.RAMUsage.Text = formatMetricValue(ramUsage)
	canvas.Refresh(m.RAMUsage)

	m.HDDTemp.Text = formatMetricValue(hddTemp)
	canvas.Refresh(m.HDDTemp)
}

// formatMetricValue форматирует метрику: если отрицательное значение, то "--"
func formatMetricValue(val float64) string {
	if val < 0 {
		return "--"
	}
	return fmt.Sprintf("%.0f", val)
}
