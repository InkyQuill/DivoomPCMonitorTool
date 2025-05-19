package metrics

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"divoom/logger"
)

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

// GetCPUInfo получает информацию о процессоре
func GetCPUInfo() string {
	lines, err := WMIQuery("cimv2", "Win32_Processor")
	if err != nil {
		logger.Log("Failed to get CPU info: %v", err)
		return "Unknown CPU"
	}

	var name, manufacturer string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Name=") {
			name = strings.TrimPrefix(line, "Name=")
		} else if strings.HasPrefix(line, "Manufacturer=") {
			manufacturer = strings.TrimPrefix(line, "Manufacturer=")
		}
	}

	if name != "" {
		logger.Log("Found CPU: %s", name)
		return name
	}
	if manufacturer != "" {
		logger.Log("Found CPU manufacturer: %s", manufacturer)
		return manufacturer + " CPU"
	}
	logger.Log("No CPU info found")
	return "Unknown CPU"
}

// GetGPUInfo получает информацию о видеокарте
func GetGPUInfo() string {
	// Сначала пробуем через nvidia-smi, так как он дает более точную информацию
	logger.Log("Trying to get GPU info via nvidia-smi")
	out, err := exec.Command("nvidia-smi", "--query-gpu=name,memory.total", "--format=csv,noheader").Output()
	if err == nil {
		info := strings.TrimSpace(string(out))
		logger.Log("Found GPU via nvidia-smi: %s", info)
		return info
	}

	// Если nvidia-smi не сработал, пробуем через WMI
	logger.Log("Trying to get GPU info via WMI")
	lines, err := WMIQuery("cimv2", "Win32_VideoController")
	if err != nil {
		logger.Log("Failed to get GPU info via WMI: %v", err)
		return "Unknown GPU"
	}

	var name, adapterRAM string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Name=") {
			name = strings.TrimPrefix(line, "Name=")
		} else if strings.HasPrefix(line, "AdapterRAM=") {
			adapterRAM = strings.TrimPrefix(line, "AdapterRAM=")
		}
	}

	if name != "" {
		// Если есть информация о памяти, добавляем её
		if adapterRAM != "" {
			if ram, err := strconv.ParseInt(adapterRAM, 10, 64); err == nil {
				// WMI возвращает размер в байтах, конвертируем в ГБ
				ramGB := float64(ram) / (1024 * 1024 * 1024)
				logger.Log("Found GPU: %s with %d bytes of memory (%.1f GB)", name, ram, ramGB)
				return fmt.Sprintf("%s (%.1f GB)", name, ramGB)
			}
		}
		logger.Log("Found GPU: %s (no memory info)", name)
		return name
	}

	logger.Log("No GPU info found")
	return "Unknown GPU"
}
