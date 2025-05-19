package metrics

import (
	"divoom/logger"
	"encoding/xml"
	"os/exec"
	"strconv"
	"strings"
)

// TemperatureMonitor интерфейс для получения температуры
type TemperatureMonitor interface {
	GetCPUTemperature() float64
	GetHDDTemperature() float64
}

// WMIMonitor реализует получение температуры через WMI
type WMIMonitor struct{}

func (w *WMIMonitor) GetCPUTemperature() float64 {
	return GetCPUTemperatureWMI()
}

func (w *WMIMonitor) GetHDDTemperature() float64 {
	return GetHDDTemperatureWMI()
}

// OpenHardwareMonitor реализует получение температуры через OpenHardwareMonitor
type OpenHardwareMonitor struct {
	path string
}

func NewOpenHardwareMonitor() *OpenHardwareMonitor {
	// Проверяем стандартные пути установки
	paths := []string{
		"C:\\Program Files\\OpenHardwareMonitor\\OpenHardwareMonitor.exe",
		"C:\\Program Files (x86)\\OpenHardwareMonitor\\OpenHardwareMonitor.exe",
	}

	for _, path := range paths {
		if _, err := exec.Command("cmd", "/C", "if exist "+path+" echo 1").Output(); err == nil {
			return &OpenHardwareMonitor{path: path}
		}
	}
	return nil
}

func (o *OpenHardwareMonitor) GetCPUTemperature() float64 {
	if o == nil {
		return -1
	}

	// Запускаем OHM и получаем XML с данными
	out, err := exec.Command(o.path, "/sensors").Output()
	if err != nil {
		return -1
	}

	// Парсим XML
	type Sensor struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	}
	type Hardware struct {
		Type    string   `xml:"type,attr"`
		Sensors []Sensor `xml:"sensor"`
	}
	type Computer struct {
		Hardware []Hardware `xml:"hardware"`
	}

	var computer Computer
	if err := xml.Unmarshal(out, &computer); err != nil {
		return -1
	}

	// Ищем температуру CPU
	for _, hw := range computer.Hardware {
		if hw.Type == "CPU" {
			for _, sensor := range hw.Sensors {
				if strings.Contains(sensor.Name, "Temperature") {
					temp, err := strconv.ParseFloat(sensor.Value, 64)
					if err == nil {
						return temp
					}
				}
			}
		}
	}
	return -1
}

func (o *OpenHardwareMonitor) GetHDDTemperature() float64 {
	if o == nil {
		return -1
	}

	out, err := exec.Command(o.path, "/sensors").Output()
	if err != nil {
		return -1
	}

	type Sensor struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	}
	type Hardware struct {
		Type    string   `xml:"type,attr"`
		Sensors []Sensor `xml:"sensor"`
	}
	type Computer struct {
		Hardware []Hardware `xml:"hardware"`
	}

	var computer Computer
	if err := xml.Unmarshal(out, &computer); err != nil {
		return -1
	}

	// Ищем температуру HDD
	for _, hw := range computer.Hardware {
		if hw.Type == "HDD" {
			for _, sensor := range hw.Sensors {
				if strings.Contains(sensor.Name, "Temperature") {
					temp, err := strconv.ParseFloat(sensor.Value, 64)
					if err == nil {
						return temp
					}
				}
			}
		}
	}
	return -1
}

// CrystalDiskInfo реализует получение температуры через CrystalDiskInfo
type CrystalDiskInfo struct {
	path string
}

func NewCrystalDiskInfo() *CrystalDiskInfo {
	// Проверяем стандартные пути установки
	paths := []string{
		"C:\\Program Files\\CrystalDiskInfo\\CrystalDiskInfo.exe",
		"C:\\Program Files (x86)\\CrystalDiskInfo\\CrystalDiskInfo.exe",
	}

	for _, path := range paths {
		if _, err := exec.Command("cmd", "/C", "if exist "+path+" echo 1").Output(); err == nil {
			return &CrystalDiskInfo{path: path}
		}
	}
	return nil
}

func (c *CrystalDiskInfo) GetCPUTemperature() float64 {
	return -1 // CrystalDiskInfo не предоставляет температуру CPU
}

func (c *CrystalDiskInfo) GetHDDTemperature() float64 {
	if c == nil {
		return -1
	}

	// Запускаем CrystalDiskInfo и получаем данные
	out, err := exec.Command(c.path, "/CopyAll").Output()
	if err != nil {
		return -1
	}

	// Ищем строку с температурой
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Temperature") {
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "Temperature" && i+1 < len(fields) {
					temp, err := strconv.ParseFloat(fields[i+1], 64)
					if err == nil {
						return temp
					}
				}
			}
		}
	}
	return -1
}

// HWiNFOMonitor реализует получение температуры через HWiNFO
type HWiNFOMonitor struct {
	path string
}

func NewHWiNFOMonitor() *HWiNFOMonitor {
	// Проверяем стандартные пути установки
	paths := []string{
		"C:\\Program Files\\HWiNFO64\\HWiNFO64.exe",
		"C:\\Program Files (x86)\\HWiNFO64\\HWiNFO64.exe",
	}

	for _, path := range paths {
		if _, err := exec.Command("cmd", "/C", "if exist "+path+" echo 1").Output(); err == nil {
			return &HWiNFOMonitor{path: path}
		}
	}
	return nil
}

func (h *HWiNFOMonitor) GetCPUTemperature() float64 {
	if h == nil {
		return -1
	}

	// Запускаем HWiNFO в режиме датчиков
	out, err := exec.Command(h.path, "/Sensors").Output()
	if err != nil {
		logger.Log("Failed to get CPU temperature from HWiNFO: %v", err)
		return -1
	}

	// Ищем строку с температурой CPU
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "CPU") && strings.Contains(line, "Temperature") {
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "Temperature" && i+1 < len(fields) {
					temp, err := strconv.ParseFloat(fields[i+1], 64)
					if err == nil {
						logger.Log("Found CPU temperature via HWiNFO: %.1f°C", temp)
						return temp
					}
				}
			}
		}
	}
	return -1
}

func (h *HWiNFOMonitor) GetHDDTemperature() float64 {
	if h == nil {
		return -1
	}

	out, err := exec.Command(h.path, "/Sensors").Output()
	if err != nil {
		logger.Log("Failed to get HDD temperature from HWiNFO: %v", err)
		return -1
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "HDD") && strings.Contains(line, "Temperature") {
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "Temperature" && i+1 < len(fields) {
					temp, err := strconv.ParseFloat(fields[i+1], 64)
					if err == nil {
						logger.Log("Found HDD temperature via HWiNFO: %.1f°C", temp)
						return temp
					}
				}
			}
		}
	}
	return -1
}

// MultiMonitor объединяет несколько мониторов
type MultiMonitor struct {
	monitors []TemperatureMonitor
}

func NewMultiMonitor() *MultiMonitor {
	monitors := []TemperatureMonitor{
		&WMIMonitor{},
		NewOpenHardwareMonitor(),
		NewCrystalDiskInfo(),
		NewHWiNFOMonitor(),
	}
	return &MultiMonitor{monitors: monitors}
}

func (m *MultiMonitor) GetCPUTemperature() float64 {
	for _, monitor := range m.monitors {
		if temp := monitor.GetCPUTemperature(); temp > 0 {
			return temp
		}
	}
	return -1
}

func (m *MultiMonitor) GetHDDTemperature() float64 {
	for _, monitor := range m.monitors {
		if temp := monitor.GetHDDTemperature(); temp > 0 {
			return temp
		}
	}
	return -1
}
