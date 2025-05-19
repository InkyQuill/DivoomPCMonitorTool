package metrics

import "time"

// SystemMetrics содержит все системные метрики
type SystemMetrics struct {
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
	// Update обновляет все метрики
	Update() error
	// GetMetrics возвращает текущие метрики
	GetMetrics() *SystemMetrics
}

// NewSystemMetrics создает новый объект метрик
func NewSystemMetrics() *SystemMetrics {
	return &SystemMetrics{
		LastUpdated: time.Now(),
	}
}
