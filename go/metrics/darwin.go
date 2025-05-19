package metrics

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// DarwinCollector реализует Collector для macOS
type DarwinCollector struct {
	metrics *SystemMetrics
}

// NewDarwinCollector создает новый коллектор для macOS
func NewDarwinCollector() *DarwinCollector {
	return &DarwinCollector{
		metrics: NewSystemMetrics(),
	}
}

// Update обновляет все метрики
func (c *DarwinCollector) Update() error {
	// CPU Usage
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercent) > 0 {
		c.metrics.CPUUsage = cpuPercent[0]
	}

	// RAM Usage
	if vmStat, err := mem.VirtualMemory(); err == nil {
		c.metrics.RAMUsage = vmStat.UsedPercent
	}

	// CPU Temperature
	c.metrics.CPUTemp = c.getCPUTemp()

	// GPU metrics
	c.metrics.GPUTemp, c.metrics.GPUUsage = c.getGPUMetrics()

	// HDD Temperature
	c.metrics.HDDTemp = c.getHDDTemp()

	c.metrics.LastUpdated = time.Now()
	return nil
}

// GetMetrics возвращает текущие метрики
func (c *DarwinCollector) GetMetrics() *SystemMetrics {
	return c.metrics
}

// getCPUTemp получает температуру CPU
func (c *DarwinCollector) getCPUTemp() float64 {
	// Пробуем получить через osx-cpu-temp
	temp := c.getCPUTempFromOSXCPUTemp()
	if temp > 0 {
		return temp
	}

	// Пробуем получить через powermetrics
	temp = c.getCPUTempFromPowerMetrics()
	if temp > 0 {
		return temp
	}

	return 0.0
}

// getCPUTempFromOSXCPUTemp получает температуру CPU через osx-cpu-temp
func (c *DarwinCollector) getCPUTempFromOSXCPUTemp() float64 {
	out, err := exec.Command("osx-cpu-temp").Output()
	if err != nil {
		return 0.0
	}

	tempStr := strings.TrimSpace(string(out))
	tempStr = strings.TrimSuffix(tempStr, "°C")
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0.0
	}

	return temp
}

// getCPUTempFromPowerMetrics получает температуру CPU через powermetrics
func (c *DarwinCollector) getCPUTempFromPowerMetrics() float64 {
	out, err := exec.Command("powermetrics", "-s", "thermal").Output()
	if err != nil {
		return 0.0
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "CPU die temperature") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				tempStr := strings.TrimSuffix(fields[4], "C")
				temp, err := strconv.ParseFloat(tempStr, 64)
				if err == nil {
					return temp
				}
			}
		}
	}

	return 0.0
}

// getGPUMetrics получает метрики GPU
func (c *DarwinCollector) getGPUMetrics() (float64, float64) {
	// Пробуем получить через powermetrics
	temp, usage := c.getGPUMetricsFromPowerMetrics()
	if temp > 0 || usage > 0 {
		return temp, usage
	}

	return 0.0, 0.0
}

// getGPUMetricsFromPowerMetrics получает метрики GPU через powermetrics
func (c *DarwinCollector) getGPUMetricsFromPowerMetrics() (float64, float64) {
	out, err := exec.Command("powermetrics", "-s", "gpu_power").Output()
	if err != nil {
		return 0.0, 0.0
	}

	var temp, usage float64
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "GPU die temperature") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				tempStr := strings.TrimSuffix(fields[4], "C")
				temp, _ = strconv.ParseFloat(tempStr, 64)
			}
		} else if strings.Contains(line, "GPU utilization") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				usageStr := strings.TrimSuffix(fields[2], "%")
				usage, _ = strconv.ParseFloat(usageStr, 64)
			}
		}
	}

	return temp, usage
}

// getHDDTemp получает температуру HDD
func (c *DarwinCollector) getHDDTemp() float64 {
	// Пробуем получить через smartctl
	temp := c.getHDDTempFromSmartctl()
	if temp > 0 {
		return temp
	}

	return 0.0
}

// getHDDTempFromSmartctl получает температуру HDD через smartctl
func (c *DarwinCollector) getHDDTempFromSmartctl() float64 {
	out, err := exec.Command("smartctl", "-A", "/dev/disk0").Output()
	if err != nil {
		return 0.0
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Temperature_Celsius") {
			fields := strings.Fields(line)
			if len(fields) >= 10 {
				temp, err := strconv.ParseFloat(fields[9], 64)
				if err == nil {
					return temp
				}
			}
		}
	}

	return 0.0
}

// GetDarwinSystemInfo возвращает редакцию MacOS
func GetDarwinSystemInfo() (edition string) {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err == nil {
		version := strings.TrimSpace(string(out))
		// Можно добавить маппинг версий к названиям (например, 14.x - Sonoma, 15.x - Sequoia)
		edition = "macOS " + version
	} else {
		edition = "macOS"
	}
	return
}
