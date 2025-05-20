package metrics

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// DarwinCollector реализует сбор информации о системе для macOS
type DarwinCollector struct {
	monitor *DarwinMonitor
	metrics *SystemMetrics
}

// NewDarwinCollector создает новый экземпляр DarwinCollector
func NewDarwinCollector() *DarwinCollector {
	return &DarwinCollector{
		monitor: &DarwinMonitor{},
		metrics: NewSystemMetrics(),
	}
}

// GetCPUInfo возвращает информацию о CPU
func (c *DarwinCollector) GetCPUInfo() string {
	return GetCPUInfoDarwin()
}

// GetGPUInfo возвращает информацию о GPU
func (c *DarwinCollector) GetGPUInfo() string {
	return GetGPUInfoDarwin()
}

// GetCPUTemperature возвращает температуру CPU
func (c *DarwinCollector) GetCPUTemperature() float64 {
	return c.monitor.GetCPUTemperature()
}

// GetHDDTemperature возвращает температуру HDD
func (c *DarwinCollector) GetHDDTemperature() float64 {
	return c.monitor.GetHDDTemperature()
}

// Update обновляет все метрики
func (c *DarwinCollector) Update() error {
	// CPU Info
	c.metrics.CPUInfo = c.GetCPUInfo()

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
	c.metrics.CPUTemp = c.GetCPUTemperature()

	// GPU metrics
	c.metrics.GPUTemp, c.metrics.GPUUsage = c.getGPUMetrics()

	// HDD Temperature
	c.metrics.HDDTemp = c.GetHDDTemperature()

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

	// Пробуем получить через nvidia-smi
	temp, usage = c.getNvidiaMetrics()
	if temp > 0 || usage > 0 {
		return temp, usage
	}

	// Пробуем получить через AMD GPU
	temp, usage = c.getAMDMetrics()
	if temp > 0 || usage > 0 {
		return temp, usage
	}

	// Пробуем получить через Intel GPU
	temp, usage = c.getIntelMetrics()
	if temp > 0 || usage > 0 {
		return temp, usage
	}

	return -1, -1
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

// getNvidiaMetrics получает метрики NVIDIA GPU
func (c *DarwinCollector) getNvidiaMetrics() (float64, float64) {
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

// getAMDMetrics получает метрики AMD GPU
func (c *DarwinCollector) getAMDMetrics() (float64, float64) {
	// В macOS для AMD GPU используем powermetrics
	out, err := exec.Command("powermetrics", "-s", "gpu_power").Output()
	if err != nil {
		return -1, -1
	}

	var temp, usage float64 = -1, -1
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

// getIntelMetrics получает метрики Intel GPU
func (c *DarwinCollector) getIntelMetrics() (float64, float64) {
	// В macOS для Intel GPU используем powermetrics
	out, err := exec.Command("powermetrics", "-s", "gpu_power").Output()
	if err != nil {
		return -1, -1
	}

	var temp, usage float64 = -1, -1
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

// GetCPUInfoDarwin получает информацию о процессоре в macOS
func GetCPUInfoDarwin() string {
	// Пробуем получить через sysctl
	out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	// Если не получилось, пробуем через system_profiler
	out, err = exec.Command("system_profiler", "SPHardwareDataType").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Processor Name:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return "Unknown Processor"
}

// GetGPUInfoDarwin получает информацию о видеокарте в macOS
func GetGPUInfoDarwin() string {
	// Пробуем получить через nvidia-smi
	out, err := exec.Command("nvidia-smi", "--query-gpu=name,memory.total", "--format=csv,noheader").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	// Пробуем получить через system_profiler
	out, err = exec.Command("system_profiler", "SPDisplaysDataType").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		var gpuInfo strings.Builder
		var currentGPU string
		var currentVRAM string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "Chipset Model:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					currentGPU = strings.TrimSpace(parts[1])
				}
			} else if strings.Contains(line, "VRAM") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					currentVRAM = strings.TrimSpace(parts[1])
				}
			} else if strings.Contains(line, "Display Type:") {
				// Это начало нового дисплея, сохраняем предыдущую информацию
				if currentGPU != "" {
					if gpuInfo.Len() > 0 {
						gpuInfo.WriteString(", ")
					}
					gpuInfo.WriteString(currentGPU)
					if currentVRAM != "" {
						gpuInfo.WriteString(" (")
						gpuInfo.WriteString(currentVRAM)
						gpuInfo.WriteString(")")
					}
					currentGPU = ""
					currentVRAM = ""
				}
			}
		}

		// Добавляем последнюю видеокарту
		if currentGPU != "" {
			if gpuInfo.Len() > 0 {
				gpuInfo.WriteString(", ")
			}
			gpuInfo.WriteString(currentGPU)
			if currentVRAM != "" {
				gpuInfo.WriteString(" (")
				gpuInfo.WriteString(currentVRAM)
				gpuInfo.WriteString(")")
			}
		}

		if gpuInfo.Len() > 0 {
			return gpuInfo.String()
		}
	}

	// Пробуем получить через powermetrics
	out, err = exec.Command("powermetrics", "-s", "gpu_power").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "GPU name:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return "Unknown GPU"
}

// DarwinMonitor реализует получение температуры в macOS
type DarwinMonitor struct{}

func (d *DarwinMonitor) GetCPUTemperature() float64 {
	// Пробуем получить через osx-cpu-temp
	temp := getCPUTempFromOSXCPUTemp()
	if temp > 0 {
		return temp
	}

	// Пробуем получить через powermetrics
	temp = getCPUTempFromPowerMetrics()
	if temp > 0 {
		return temp
	}

	return -1
}

func (d *DarwinMonitor) GetHDDTemperature() float64 {
	// Пробуем получить через smartctl
	temp := GetHDDTempFromSmartctl("/dev/disk0")
	if temp > 0 {
		return temp
	}

	return -1
}

// getCPUTempFromOSXCPUTemp получает температуру CPU через osx-cpu-temp
func getCPUTempFromOSXCPUTemp() float64 {
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
func getCPUTempFromPowerMetrics() float64 {
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
