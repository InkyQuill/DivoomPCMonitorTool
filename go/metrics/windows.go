package metrics

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// WindowsCollector реализует Collector для Windows
type WindowsCollector struct {
	metrics *SystemMetrics
	monitor *MultiMonitor
}

// NewWindowsCollector создает новый коллектор для Windows
func NewWindowsCollector() *WindowsCollector {
	return &WindowsCollector{
		metrics: NewSystemMetrics(),
		monitor: NewMultiMonitor(),
	}
}

// Update обновляет все метрики
func (c *WindowsCollector) Update() error {
	// CPU Usage
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercent) > 0 {
		c.metrics.CPUUsage = cpuPercent[0]
	} else {
		c.metrics.CPUUsage = -1
	}

	// RAM Usage
	if vmStat, err := mem.VirtualMemory(); err == nil {
		c.metrics.RAMUsage = vmStat.UsedPercent
	} else {
		c.metrics.RAMUsage = -1
	}

	// CPU Temperature
	c.metrics.CPUTemp = c.monitor.GetCPUTemperature()

	// GPU metrics
	c.metrics.GPUTemp, c.metrics.GPUUsage = c.getGPUMetrics()

	// HDD Temperature
	c.metrics.HDDTemp = c.monitor.GetHDDTemperature()

	c.metrics.LastUpdated = time.Now()
	return nil
}

// GetMetrics возвращает текущие метрики
func (c *WindowsCollector) GetMetrics() *SystemMetrics {
	return c.metrics
}

// getCPUTemp получает температуру CPU через WMI
func (c *WindowsCollector) getCPUTemp() float64 {
	// Пробуем получить через WMI
	out, err := exec.Command("wmic", "/namespace:\\\\root\\wmi", "path", "MSAcpi_ThermalZoneTemperature", "get", "CurrentTemperature").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" || strings.Contains(line, "CurrentTemperature") {
				continue
			}
			temp, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
			if err == nil {
				// Конвертируем из десятых градуса Кельвина в градусы Цельсия
				return (temp/10.0 - 273.15)
			}
		}
	}

	// Если WMI не сработал, пробуем через OpenHardwareMonitor
	// TODO: Добавить поддержку OpenHardwareMonitor
	return -1
}

// getGPUMetrics получает метрики GPU
func (c *WindowsCollector) getGPUMetrics() (float64, float64) {
	// Пробуем получить через nvidia-smi
	temp, usage := c.getNvidiaMetrics()
	if temp > 0 || usage > 0 {
		return temp, usage
	}

	// TODO: Добавить поддержку AMD GPU через Radeon Software
	// TODO: Добавить поддержку Intel GPU через Intel Graphics Command Center

	return -1, -1
}

// getNvidiaMetrics получает метрики NVIDIA GPU
func (c *WindowsCollector) getNvidiaMetrics() (float64, float64) {
	// Температура
	tempOut, err := exec.Command("nvidia-smi", "--query-gpu=temperature.gpu", "--format=csv,noheader,nounits").Output()
	if err != nil {
		return -1, -1
	}
	tempStr := strings.TrimSpace(string(tempOut))
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		temp = -1
	}

	// Загрузка
	usageOut, err := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits").Output()
	if err != nil {
		return temp, -1
	}
	usageStr := strings.TrimSpace(string(usageOut))
	usage, err := strconv.ParseFloat(usageStr, 64)
	if err != nil {
		usage = -1
	}

	return temp, usage
}

// getHDDTemp получает температуру HDD через WMI
func (c *WindowsCollector) getHDDTemp() float64 {
	// Пробуем получить через WMI
	out, err := exec.Command("wmic", "diskdrive", "get", "Temperature").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" || strings.Contains(line, "Temperature") {
				continue
			}
			temp, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
			if err == nil {
				return temp
			}
		}
	}

	// Если WMI не сработал, пробуем через CrystalDiskInfo
	// TODO: Добавить поддержку CrystalDiskInfo
	return -1
}

// GetWindowsSystemInfo возвращает редакцию, билд и архитектуру Windows
func GetWindowsSystemInfo() (edition, build, arch string) {
	// Получаем редакцию
	out, err := exec.Command("wmic", "os", "get", "Caption", "/value").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Caption=") {
				edition = strings.TrimPrefix(line, "Caption=")
				edition = strings.TrimSpace(edition)
				break
			}
		}
	}
	if edition == "" {
		edition = "Windows"
	}

	// Получаем версию
	out, err = exec.Command("wmic", "os", "get", "Version", "/value").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Version=") {
				build = strings.TrimPrefix(line, "Version=")
				build = strings.TrimSpace(build)
				break
			}
		}
	}
	if build == "" {
		build = "Unknown"
	}

	// Получаем архитектуру
	out, err = exec.Command("wmic", "os", "get", "OSArchitecture", "/value").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "OSArchitecture=") {
				arch = strings.TrimPrefix(line, "OSArchitecture=")
				arch = strings.TrimSpace(arch)
				break
			}
		}
	}
	if arch == "" {
		arch = "Unknown"
	}
	return
}
