package metrics

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// LinuxCollector реализует Collector для Linux
type LinuxCollector struct {
	metrics *SystemMetrics
}

// NewLinuxCollector создает новый коллектор для Linux
func NewLinuxCollector() *LinuxCollector {
	return &LinuxCollector{
		metrics: NewSystemMetrics(),
	}
}

// Update обновляет все метрики
func (c *LinuxCollector) Update() error {
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
	out, err := exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp").Output()
	if err != nil {
		return 0.0
	}

	tempStr := strings.TrimSpace(string(out))
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0.0
	}

	// Температура в миллиградусах
	return temp / 1000.0
}

// getGPUMetrics получает метрики GPU
func (c *LinuxCollector) getGPUMetrics() (float64, float64) {
	// Пробуем получить через nvidia-smi
	temp, usage := c.getNvidiaMetrics()
	if temp > 0 || usage > 0 {
		return temp, usage
	}

	// TODO: Добавить поддержку AMD GPU через rocm-smi
	// TODO: Добавить поддержку Intel GPU через intel_gpu_top

	return 0.0, 0.0
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

// getHDDTemp получает температуру HDD
func (c *LinuxCollector) getHDDTemp() float64 {
	// Пробуем получить через hddtemp
	temp := c.getHDDTempFromHDDTemp()
	if temp > 0 {
		return temp
	}

	// Пробуем получить через smartctl
	temp = c.getHDDTempFromSmartctl()
	if temp > 0 {
		return temp
	}

	return 0.0
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
func (c *LinuxCollector) getHDDTempFromSmartctl() float64 {
	out, err := exec.Command("smartctl", "-A", "/dev/sda").Output()
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
