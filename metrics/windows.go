package metrics

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"divoom/logger"
)

// WindowsCollector реализует сбор информации о системе для Windows
type WindowsCollector struct {
	monitor *WindowsMonitor
	metrics *SystemMetrics
}

// NewWindowsCollector создает новый экземпляр WindowsCollector
func NewWindowsCollector() *WindowsCollector {
	return &WindowsCollector{
		monitor: &WindowsMonitor{},
		metrics: NewSystemMetrics(),
	}
}

// WMIQuery выполняет запрос к WMI и возвращает результат
func WMIQuery(namespace, query string) ([]string, error) {
	// Используем 64-битный wmic
	cmd := exec.Command("C:\\Windows\\System32\\wbem\\wmic.exe", "/namespace:\\\\root\\"+namespace, "path", query, "get", "/value")
	out, err := cmd.Output()
	if err != nil {
		logger.Log("WMI query failed for %s/%s: %v", namespace, query, err)
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	logger.Log("WMI query %s/%s returned %d lines", namespace, query, len(lines))
	return lines, nil
}

// GetCPUTemperatureWMI получает температуру CPU через WMI
func GetCPUTemperatureWMI() float64 {
	// Пробуем разные пространства имен и запросы
	queries := []struct {
		namespace string
		query     string
	}{
		{"wmi", "MSAcpi_ThermalZoneTemperature"},
		{"cimv2", "Win32_TemperatureProbe"},
		{"cimv2", "Win32_PerfFormattedData_Counters_ThermalZoneInformation"},
		{"cimv2", "Win32_PerfRawData_Counters_ThermalZoneInformation"},
	}

	for _, q := range queries {
		logger.Log("Trying CPU temperature query: %s/%s", q.namespace, q.query)
		lines, err := WMIQuery(q.namespace, q.query)
		if err != nil {
			logger.Log("Query failed: %v", err)
			continue
		}

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			logger.Log("Processing line: %s", line)

			// Пробуем разные форматы данных
			if strings.HasPrefix(line, "CurrentTemperature=") {
				temp, err := strconv.ParseFloat(strings.TrimPrefix(line, "CurrentTemperature="), 64)
				if err == nil {
					result := (temp/10.0 - 273.15) // Конвертируем из десятых градуса Кельвина в Цельсий
					logger.Log("Found CPU temperature (Kelvin): %.2f°C", result)
					return result
				}
			} else if strings.HasPrefix(line, "Temperature=") {
				temp, err := strconv.ParseFloat(strings.TrimPrefix(line, "Temperature="), 64)
				if err == nil {
					logger.Log("Found CPU temperature: %.2f°C", temp)
					return temp
				}
			}
		}
	}
	logger.Log("No CPU temperature data found")
	return -1
}

// GetHDDTemperatureWMI получает температуру HDD через WMI
func GetHDDTemperatureWMI() float64 {
	lines, err := WMIQuery("cimv2", "Win32_DiskDrive")
	if err != nil {
		logger.Log("Failed to get HDD temperature: %v", err)
		return -1
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Temperature=") {
			temp, err := strconv.ParseFloat(strings.TrimPrefix(line, "Temperature="), 64)
			if err == nil {
				logger.Log("Found HDD temperature: %.2f°C", temp)
				return temp
			}
		}
	}
	logger.Log("No HDD temperature data found")
	return -1
}

// GetCPUInfo возвращает информацию о CPU
func (c *WindowsCollector) GetCPUInfo() string {
	return GetCPUInfoWindows()
}

// GetGPUInfo возвращает информацию о GPU
func (c *WindowsCollector) GetGPUInfo() string {
	return GetGPUInfoWindows()
}

// GetCPUTemperature возвращает температуру CPU
func (c *WindowsCollector) GetCPUTemperature() float64 {
	return c.monitor.GetCPUTemperature()
}

// GetHDDTemperature возвращает температуру HDD
func (c *WindowsCollector) GetHDDTemperature() float64 {
	return c.monitor.GetHDDTemperature()
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

// GetCPUInfoWindows получает информацию о процессоре в Windows
func GetCPUInfoWindows() string {
	// Пробуем получить через WMI
	out, err := exec.Command("wmic", "cpu", "get", "Name", "/value").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Name=") {
				return strings.TrimSpace(strings.TrimPrefix(line, "Name="))
			}
		}
	}

	// Если не получилось, пробуем через PowerShell
	out, err = exec.Command("powershell", "-Command", "Get-WmiObject -Class Win32_Processor | Select-Object -ExpandProperty Name").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	return "Unknown Processor"
}

// GetGPUInfoWindows получает информацию о видеокарте в Windows
func GetGPUInfoWindows() string {
	// Пробуем получить через nvidia-smi
	out, err := exec.Command("nvidia-smi", "--query-gpu=name,memory.total", "--format=csv,noheader").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	// Пробуем получить через WMI
	lines, err := WMIQuery("cimv2", "Win32_VideoController")
	if err == nil {
		var gpuInfo strings.Builder
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Name=") {
				gpuInfo.WriteString(strings.TrimPrefix(line, "Name="))
			} else if strings.HasPrefix(line, "AdapterRAM=") {
				ram := strings.TrimPrefix(line, "AdapterRAM=")
				if ram != "" {
					gpuInfo.WriteString(" (")
					gpuInfo.WriteString(ram)
					gpuInfo.WriteString(")")
				}
			}
		}
		if gpuInfo.Len() > 0 {
			return gpuInfo.String()
		}
	}

	// Пробуем получить через PowerShell
	out, err = exec.Command("powershell", "-Command", "Get-WmiObject -Class Win32_VideoController | Select-Object Name, AdapterRAM | Format-List").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		var gpuInfo strings.Builder
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Name") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					gpuInfo.WriteString(strings.TrimSpace(parts[1]))
				}
			} else if strings.HasPrefix(line, "AdapterRAM") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					ram := strings.TrimSpace(parts[1])
					if ram != "" {
						gpuInfo.WriteString(" (")
						gpuInfo.WriteString(ram)
						gpuInfo.WriteString(")")
					}
				}
			}
		}
		if gpuInfo.Len() > 0 {
			return gpuInfo.String()
		}
	}

	return "Unknown GPU"
}

// WindowsMonitor реализует получение температуры в Windows
type WindowsMonitor struct{}

func (w *WindowsMonitor) GetCPUTemperature() float64 {
	// Пробуем получить через WMI
	temp := GetCPUTemperatureWMI()
	if temp > 0 {
		return temp
	}

	// Пробуем получить через OpenHardwareMonitor
	temp = getCPUTempFromOpenHardwareMonitor()
	if temp > 0 {
		return temp
	}

	return -1
}

func (w *WindowsMonitor) GetHDDTemperature() float64 {
	// Пробуем получить через WMI
	temp := GetHDDTemperatureWMI()
	if temp > 0 {
		return temp
	}

	// Пробуем получить через OpenHardwareMonitor
	temp = getHDDTempFromOpenHardwareMonitor()
	if temp > 0 {
		return temp
	}

	return -1
}

// getCPUTempFromOpenHardwareMonitor получает температуру CPU через OpenHardwareMonitor
func getCPUTempFromOpenHardwareMonitor() float64 {
	// TODO: Реализовать получение температуры через OpenHardwareMonitor
	return 0.0
}

// getHDDTempFromOpenHardwareMonitor получает температуру HDD через OpenHardwareMonitor
func getHDDTempFromOpenHardwareMonitor() float64 {
	// TODO: Реализовать получение температуры через OpenHardwareMonitor
	return 0.0
}

// Update обновляет все метрики
func (c *WindowsCollector) Update() error {
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
func (c *WindowsCollector) GetMetrics() *SystemMetrics {
	return c.metrics
}

// getGPUMetrics получает метрики GPU
func (c *WindowsCollector) getGPUMetrics() (float64, float64) {
	// Пробуем получить через nvidia-smi
	temp, usage := c.getNvidiaMetrics()
	if temp > 0 || usage > 0 {
		return temp, usage
	}

	// Пробуем получить через Radeon Software для AMD
	temp, usage = c.getAMDMetrics()
	if temp > 0 || usage > 0 {
		return temp, usage
	}

	// Пробуем получить через Intel Graphics Command Center
	temp, usage = c.getIntelMetrics()
	if temp > 0 || usage > 0 {
		return temp, usage
	}

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

// getAMDMetrics получает метрики AMD GPU
func (c *WindowsCollector) getAMDMetrics() (float64, float64) {
	// Пробуем получить через WMI
	lines, err := WMIQuery("cimv2", "Win32_VideoController")
	if err != nil {
		return -1, -1
	}

	var temp, usage float64 = -1, -1
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Temperature=") {
			temp, _ = strconv.ParseFloat(strings.TrimPrefix(line, "Temperature="), 64)
		} else if strings.HasPrefix(line, "Utilization=") {
			usage, _ = strconv.ParseFloat(strings.TrimPrefix(line, "Utilization="), 64)
		}
	}

	// Если не получилось через WMI, пробуем через Radeon Software
	if temp < 0 || usage < 0 {
		// TODO: Добавить поддержку Radeon Software API
		// Это потребует использования COM-интерфейсов или чтения файлов конфигурации
	}

	return temp, usage
}

// getIntelMetrics получает метрики Intel GPU
func (c *WindowsCollector) getIntelMetrics() (float64, float64) {
	// Пробуем получить через WMI
	lines, err := WMIQuery("cimv2", "Win32_VideoController")
	if err != nil {
		return -1, -1
	}

	var temp, usage float64 = -1, -1
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Temperature=") {
			temp, _ = strconv.ParseFloat(strings.TrimPrefix(line, "Temperature="), 64)
		} else if strings.HasPrefix(line, "Utilization=") {
			usage, _ = strconv.ParseFloat(strings.TrimPrefix(line, "Utilization="), 64)
		}
	}

	// Если не получилось через WMI, пробуем через Intel Graphics Command Center
	if temp < 0 || usage < 0 {
		// TODO: Добавить поддержку Intel Graphics Command Center API
		// Это потребует использования COM-интерфейсов или чтения файлов конфигурации
	}

	return temp, usage
}
