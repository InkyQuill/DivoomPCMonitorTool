package metrics

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// LinuxCollector реализует сбор информации о системе для Linux
type LinuxCollector struct {
	monitor *LinuxMonitor
	metrics *SystemMetrics
}

// NewLinuxCollector создает новый экземпляр LinuxCollector
func NewLinuxCollector() *LinuxCollector {
	return &LinuxCollector{
		monitor: &LinuxMonitor{},
		metrics: NewSystemMetrics(),
	}
}

// GetCPUInfo возвращает информацию о CPU
func (c *LinuxCollector) GetCPUInfo() string {
	return GetCPUInfoLinux()
}

// GetGPUInfo возвращает информацию о GPU
func (c *LinuxCollector) GetGPUInfo() string {
	return GetGPUInfoLinux()
}

// GetCPUTemperature возвращает температуру CPU
func (c *LinuxCollector) GetCPUTemperature() float64 {
	return c.monitor.GetCPUTemperature()
}

// GetHDDTemperature возвращает температуру HDD
func (c *LinuxCollector) GetHDDTemperature() float64 {
	return c.monitor.GetHDDTemperature()
}

// Update обновляет все метрики
func (c *LinuxCollector) Update() error {
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
func (c *LinuxCollector) GetMetrics() *SystemMetrics {
	return c.metrics
}

// getCPUTemp получает температуру CPU
func (c *LinuxCollector) getCPUTemp() float64 {
	// Пробуем получить через sensors
	temp := c.getCPUTempFromSensors()
	if temp > 0 {
		return temp
	}

	// Пробуем получить через thermal
	temp = c.getCPUTempFromThermal()
	if temp > 0 {
		return temp
	}

	return 0.0
}

// getCPUTempFromSensors получает температуру CPU через lm-sensors
func (c *LinuxCollector) getCPUTempFromSensors() float64 {
	// Пробуем разные паттерны для разных CPU
	patterns := []string{
		"sensors | grep -m1 'Package id 0:' | awk '{print $4}' | tr -d '+°C'",
		"sensors | grep -m1 'Core 0:' | awk '{print $3}' | tr -d '+°C'",
		"sensors | grep -m1 'Tdie:' | awk '{print $2}' | tr -d '+°C'",
		"sensors | grep -m1 'CPU Temperature:' | awk '{print $3}' | tr -d '+°C'",
		"sensors | grep -m1 'CPU:' | awk '{print $2}' | tr -d '+°C'",
		"sensors | grep -m1 'temp1:' | awk '{print $2}' | tr -d '+°C'",
		"sensors | grep -m1 'Tctl:' | awk '{print $2}' | tr -d '+°C'",  // AMD
		"sensors | grep -m1 'Tccd1:' | awk '{print $2}' | tr -d '+°C'", // AMD CCD
	}

	for _, pattern := range patterns {
		out, err := exec.Command("bash", "-c", pattern).Output()
		if err == nil {
			tempStr := strings.TrimSpace(string(out))
			temp, err := strconv.ParseFloat(tempStr, 64)
			if err == nil && temp > 0 {
				return temp
			}
		}
	}

	return 0.0
}

// getCPUTempFromThermal получает температуру CPU через thermal
func (c *LinuxCollector) getCPUTempFromThermal() float64 {
	// Пробуем разные thermal zones
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("/sys/class/thermal/thermal_zone%d/temp", i)
		out, err := exec.Command("cat", path).Output()
		if err != nil {
			continue
		}

		tempStr := strings.TrimSpace(string(out))
		temp, err := strconv.ParseFloat(tempStr, 64)
		if err == nil && temp > 0 {
			// Проверяем, что это действительно CPU
			typePath := fmt.Sprintf("/sys/class/thermal/thermal_zone%d/type", i)
			typeOut, err := exec.Command("cat", typePath).Output()
			if err == nil {
				typeStr := strings.TrimSpace(string(typeOut))
				if strings.Contains(strings.ToLower(typeStr), "cpu") ||
					strings.Contains(strings.ToLower(typeStr), "x86") ||
					strings.Contains(strings.ToLower(typeStr), "acpitz") {
					return temp / 1000.0
				}
			}
		}
	}

	return 0.0
}

// getGPUMetrics получает метрики GPU
func (c *LinuxCollector) getGPUMetrics() (float64, float64) {
	// Пробуем получить через nvidia-smi
	temp, usage := c.getNvidiaMetrics()
	if temp > 0 || usage > 0 {
		return temp, usage
	}

	// Пробуем получить через rocm-smi для AMD
	temp, usage = c.getAMDMetrics()
	if temp > 0 || usage > 0 {
		return temp, usage
	}

	// Пробуем получить через intel_gpu_top для Intel
	temp, usage = c.getIntelMetrics()
	if temp > 0 || usage > 0 {
		return temp, usage
	}

	return -1, -1
}

// getNvidiaMetrics получает метрики NVIDIA GPU
func (c *LinuxCollector) getNvidiaMetrics() (float64, float64) {
	// Температура
	tempOut, err := exec.Command("nvidia-smi", "--query-gpu=temperature.gpu", "--format=csv,noheader,nounits").Output()
	if err != nil {
		return 0.0, 0.0
	}
	tempStr := strings.TrimSpace(string(tempOut))
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		temp = 0.0
	}

	// Загрузка
	usageOut, err := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits").Output()
	if err != nil {
		return temp, 0.0
	}
	usageStr := strings.TrimSpace(string(usageOut))
	usage, err := strconv.ParseFloat(usageStr, 64)
	if err != nil {
		usage = 0.0
	}

	return temp, usage
}

// getAMDMetrics получает метрики AMD GPU
func (c *LinuxCollector) getAMDMetrics() (float64, float64) {
	// Температура
	tempOut, err := exec.Command("rocm-smi", "--showtemp", "--json").Output()
	if err != nil {
		return -1, -1
	}

	// Парсим JSON для получения температуры
	var result map[string]interface{}
	if err := json.Unmarshal(tempOut, &result); err != nil {
		return -1, -1
	}

	temp := -1.0
	if gpuTemp, ok := result["GPU Temperature"].(float64); ok {
		temp = gpuTemp
	}

	// Загрузка
	usageOut, err := exec.Command("rocm-smi", "--showuse", "--json").Output()
	if err != nil {
		return temp, -1
	}

	// Парсим JSON для получения загрузки
	if err := json.Unmarshal(usageOut, &result); err != nil {
		return temp, -1
	}

	usage := -1.0
	if gpuUse, ok := result["GPU use (%)"].(float64); ok {
		usage = gpuUse
	}

	return temp, usage
}

// getIntelMetrics получает метрики Intel GPU
func (c *LinuxCollector) getIntelMetrics() (float64, float64) {
	// Температура
	tempOut, err := exec.Command("intel_gpu_top", "-J", "-s", "1").Output()
	if err != nil {
		return -1, -1
	}

	// Парсим JSON для получения температуры и загрузки
	var result map[string]interface{}
	if err := json.Unmarshal(tempOut, &result); err != nil {
		return -1, -1
	}

	temp := -1.0
	usage := -1.0

	if engines, ok := result["engines"].(map[string]interface{}); ok {
		if render, ok := engines["Render/3D/0"].(map[string]interface{}); ok {
			if busy, ok := render["busy"].(float64); ok {
				usage = busy
			}
		}
	}

	// Для температуры используем sensors
	sensorsOut, err := exec.Command("sensors").Output()
	if err == nil {
		lines := strings.Split(string(sensorsOut), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Package id 0:") {
				fields := strings.Fields(line)
				if len(fields) >= 4 {
					tempStr := strings.Trim(fields[3], "+°C")
					if t, err := strconv.ParseFloat(tempStr, 64); err == nil {
						temp = t
					}
				}
			}
		}
	}

	return temp, usage
}

// getHDDTemp получает температуру HDD
func (c *LinuxCollector) getHDDTemp() float64 {
	// Пробуем получить через NVMe
	temp := c.getNVMeTemp()
	if temp > 0 {
		return temp
	}

	// Пробуем получить через hddtemp
	temp = c.getHDDTempFromHDDTemp()
	if temp > 0 {
		return temp
	}

	// Пробуем получить через smartctl
	temp = c.getHDDTempFromSmartctl("/dev/sda")
	if temp > 0 {
		return temp
	}

	return 0.0
}

// getNVMeTemp получает температуру NVMe накопителей
func (c *LinuxCollector) getNVMeTemp() float64 {
	// Пробуем получить через sensors
	out, err := exec.Command("sensors").Output()
	if err != nil {
		return 0.0
	}

	lines := strings.Split(string(out), "\n")
	var maxTemp float64 = 0.0

	for _, line := range lines {
		if strings.Contains(line, "Composite:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				tempStr := strings.Trim(fields[1], "+°C")
				temp, err := strconv.ParseFloat(tempStr, 64)
				if err == nil && temp > maxTemp {
					maxTemp = temp
				}
			}
		}
	}

	return maxTemp
}

// getHDDTempFromHDDTemp получает температуру HDD через hddtemp
func (c *LinuxCollector) getHDDTempFromHDDTemp() float64 {
	out, err := exec.Command("hddtemp", "-n", "/dev/sda").Output()
	if err != nil {
		return 0.0
	}

	tempStr := strings.TrimSpace(string(out))
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0.0
	}

	return temp
}

// getHDDTempFromSmartctl получает температуру HDD через smartctl
func (c *LinuxCollector) getHDDTempFromSmartctl(device string) float64 {
	out, err := exec.Command("smartctl", "-A", device).Output()
	if err != nil {
		return 0.0
	}

	// Ищем строку с температурой
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

// GetLinuxSystemInfo возвращает редакцию и версию ядра Linux
func GetLinuxSystemInfo() (distro, kernel string) {
	// Получаем редакцию
	out, err := exec.Command("lsb_release", "-ds").Output()
	if err == nil {
		distro = strings.Trim(string(out), "\n\"")
	} else {
		// Альтернативный способ
		out, err = exec.Command("cat", "/etc/os-release").Output()
		if err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					distro = strings.Trim(line[13:], "\"")
					break
				}
			}
		}
	}
	if distro == "" {
		distro = "Unknown Linux"
	}

	// Получаем версию ядра
	out, err = exec.Command("uname", "-r").Output()
	if err == nil {
		kernel = strings.TrimSpace(string(out))
	} else {
		kernel = "Unknown"
	}
	return
}

// GetCPUInfoLinux получает информацию о процессоре в Linux
func GetCPUInfoLinux() string {
	// Пробуем получить через /proc/cpuinfo
	out, err := exec.Command("cat", "/proc/cpuinfo").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "model name") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	// Если не получилось, пробуем через lscpu
	out, err = exec.Command("lscpu").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Model name:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return "Unknown Processor"
}

// GetGPUInfoLinux получает информацию о видеокарте в Linux
func GetGPUInfoLinux() string {
	// Сначала пробуем nvidia-smi
	out, err := exec.Command("nvidia-smi", "--query-gpu=name,memory.total", "--format=csv,noheader").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	// Пробуем получить через lspci
	out, err = exec.Command("lspci", "-v").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		var gpuInfo strings.Builder
		for i, line := range lines {
			if strings.Contains(line, "VGA") || strings.Contains(line, "3D") {
				// Ищем производителя и модель
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					gpuInfo.WriteString(strings.TrimSpace(parts[1]))
					// Ищем дополнительную информацию в следующих строках
					for j := i + 1; j < len(lines) && !strings.Contains(lines[j], "VGA") && !strings.Contains(lines[j], "3D"); j++ {
						if strings.Contains(lines[j], "Subsystem") {
							subParts := strings.Split(lines[j], ":")
							if len(subParts) > 1 {
								gpuInfo.WriteString(" (")
								gpuInfo.WriteString(strings.TrimSpace(subParts[1]))
								gpuInfo.WriteString(")")
							}
							break
						}
					}
					return gpuInfo.String()
				}
			}
		}
	}

	// Пробуем получить через glxinfo
	out, err = exec.Command("glxinfo", "-B").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "OpenGL renderer string:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return "Unknown GPU"
}

// LinuxMonitor реализует получение температуры в Linux
type LinuxMonitor struct{}

func (l *LinuxMonitor) GetCPUTemperature() float64 {
	// Пробуем получить через sensors
	temp := getCPUTempFromSensors()
	if temp > 0 {
		return temp
	}

	// Пробуем получить через thermal
	temp = getCPUTempFromThermal()
	if temp > 0 {
		return temp
	}

	return -1
}

func (l *LinuxMonitor) GetHDDTemperature() float64 {
	// Пробуем получить через NVMe
	temp := getNVMeTemp()
	if temp > 0 {
		return temp
	}

	// Пробуем получить через hddtemp
	temp = getHDDTempFromHDDTemp()
	if temp > 0 {
		return temp
	}

	// Пробуем получить через smartctl
	temp = GetHDDTempFromSmartctl("/dev/sda")
	if temp > 0 {
		return temp
	}

	return -1
}

// getCPUTempFromSensors получает температуру CPU через lm-sensors
func getCPUTempFromSensors() float64 {
	// Пробуем разные паттерны для разных CPU
	patterns := []string{
		"sensors | grep -m1 'Package id 0:' | awk '{print $4}' | tr -d '+°C'",
		"sensors | grep -m1 'Core 0:' | awk '{print $3}' | tr -d '+°C'",
		"sensors | grep -m1 'Tdie:' | awk '{print $2}' | tr -d '+°C'",
		"sensors | grep -m1 'CPU Temperature:' | awk '{print $3}' | tr -d '+°C'",
		"sensors | grep -m1 'CPU:' | awk '{print $2}' | tr -d '+°C'",
		"sensors | grep -m1 'temp1:' | awk '{print $2}' | tr -d '+°C'",
		"sensors | grep -m1 'Tctl:' | awk '{print $2}' | tr -d '+°C'",  // AMD
		"sensors | grep -m1 'Tccd1:' | awk '{print $2}' | tr -d '+°C'", // AMD CCD
	}

	for _, pattern := range patterns {
		out, err := exec.Command("bash", "-c", pattern).Output()
		if err == nil {
			tempStr := strings.TrimSpace(string(out))
			temp, err := strconv.ParseFloat(tempStr, 64)
			if err == nil && temp > 0 {
				return temp
			}
		}
	}

	return 0.0
}

// getCPUTempFromThermal получает температуру CPU через thermal
func getCPUTempFromThermal() float64 {
	// Пробуем разные thermal zones
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("/sys/class/thermal/thermal_zone%d/temp", i)
		out, err := exec.Command("cat", path).Output()
		if err != nil {
			continue
		}

		tempStr := strings.TrimSpace(string(out))
		temp, err := strconv.ParseFloat(tempStr, 64)
		if err == nil && temp > 0 {
			// Проверяем, что это действительно CPU
			typePath := fmt.Sprintf("/sys/class/thermal/thermal_zone%d/type", i)
			typeOut, err := exec.Command("cat", typePath).Output()
			if err == nil {
				typeStr := strings.TrimSpace(string(typeOut))
				if strings.Contains(strings.ToLower(typeStr), "cpu") ||
					strings.Contains(strings.ToLower(typeStr), "x86") ||
					strings.Contains(strings.ToLower(typeStr), "acpitz") {
					return temp / 1000.0
				}
			}
		}
	}

	return 0.0
}

// getNVMeTemp получает температуру NVMe накопителей
func getNVMeTemp() float64 {
	// Пробуем получить через sensors
	out, err := exec.Command("sensors").Output()
	if err != nil {
		return 0.0
	}

	lines := strings.Split(string(out), "\n")
	var maxTemp float64 = 0.0

	for _, line := range lines {
		if strings.Contains(line, "Composite:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				tempStr := strings.Trim(fields[1], "+°C")
				temp, err := strconv.ParseFloat(tempStr, 64)
				if err == nil && temp > maxTemp {
					maxTemp = temp
				}
			}
		}
	}

	return maxTemp
}

// getHDDTempFromHDDTemp получает температуру HDD через hddtemp
func getHDDTempFromHDDTemp() float64 {
	out, err := exec.Command("hddtemp", "-n", "/dev/sda").Output()
	if err != nil {
		return 0.0
	}

	tempStr := strings.TrimSpace(string(out))
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0.0
	}

	return temp
}
