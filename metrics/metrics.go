package metrics

import (
	"runtime"
	"time"
)

// SystemMetrics содержит все системные метрики
type SystemMetrics struct {
	CPUInfo     string
	CPUTemp     float64
	CPUUsage    float64
	GPUTemp     float64
	GPUUsage    float64
	RAMUsage    float64
	HDDTemp     float64
	LastUpdated time.Time
}

// Collector интерфейс для сбора метрик
type Collector interface {
	GetCPUInfo() string
	GetGPUInfo() string
	GetCPUTemperature() float64
	GetHDDTemperature() float64
	Update() error
	GetMetrics() *SystemMetrics
}

// NewSystemMetrics создает новый объект метрик
func NewSystemMetrics() *SystemMetrics {
	return &SystemMetrics{
		LastUpdated: time.Now(),
	}
}

// GetCPUInfo возвращает информацию о процессоре для текущей ОС
func GetCPUInfo() string {
	var collector Collector
	switch runtime.GOOS {
	case "linux":
		collector = &LinuxCollector{}
	case "windows":
		collector = &WindowsCollector{}
	case "darwin":
		collector = &DarwinCollector{}
	default:
		return "Unknown CPU"
	}
	return collector.GetCPUInfo()
}

// GetGPUInfo возвращает информацию о видеокарте для текущей ОС
func GetGPUInfo() string {
	var collector Collector
	switch runtime.GOOS {
	case "linux":
		collector = &LinuxCollector{}
	case "windows":
		collector = &WindowsCollector{}
	case "darwin":
		collector = &DarwinCollector{}
	default:
		return "Unknown GPU"
	}
	return collector.GetGPUInfo()
}

// GetSystemMetrics возвращает все метрики системы
func GetSystemMetrics() SystemMetrics {
	var collector Collector
	switch runtime.GOOS {
	case "linux":
		collector = &LinuxCollector{}
	case "windows":
		collector = &WindowsCollector{}
	case "darwin":
		collector = &DarwinCollector{}
	default:
		return SystemMetrics{
			CPUInfo:     "Unknown CPU",
			CPUTemp:     -1,
			CPUUsage:    -1,
			GPUTemp:     -1,
			GPUUsage:    -1,
			RAMUsage:    -1,
			HDDTemp:     -1,
			LastUpdated: time.Now(),
		}
	}

	return SystemMetrics{
		CPUInfo:     collector.GetCPUInfo(),
		CPUTemp:     collector.GetCPUTemperature(),
		CPUUsage:    -1,
		GPUTemp:     -1,
		GPUUsage:    -1,
		RAMUsage:    -1,
		HDDTemp:     collector.GetHDDTemperature(),
		LastUpdated: time.Now(),
	}
}
