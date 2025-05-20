package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	DeviceTypeTimeGate = 400
)

// Device представляет информацию об устройстве Divoom
type Device struct {
	DeviceName      string `json:"DeviceName"`
	DevicePrivateIP string `json:"DevicePrivateIP"`
	Hardware        int    `json:"Hardware"`
}

// Manager управляет устройствами Divoom
type Manager struct {
	client *http.Client
}

// NewManager создает новый менеджер устройств
func NewManager() *Manager {
	return &Manager{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// DiscoverDevices ищет устройства Divoom в локальной сети
func (m *Manager) DiscoverDevices() ([]Device, error) {
	resp, err := m.client.Get("http://app.divoom-gz.com/Device/ReturnSameLANDevice")
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске устройств: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка сервера: %d", resp.StatusCode)
	}

	var result struct {
		DeviceList []Device `json:"DeviceList"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %v", err)
	}

	return result.DeviceList, nil
}

// SendSystemInfo отправляет системную информацию на устройство
func (m *Manager) SendSystemInfo(ip string, lcdID int, metrics map[string]float64) error {
	data := map[string]interface{}{
		"Command": "Device/UpdatePCParaInfo",
		"ScreenList": []map[string]interface{}{
			{
				"LcdId": lcdID,
				"DispData": []float64{
					metrics["cpu_usage"],
					metrics["gpu_usage"],
					metrics["cpu_temp"],
					metrics["gpu_temp"],
					metrics["memory_usage"],
					metrics["disk_usage"],
				},
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка при подготовке данных: %v", err)
	}

	resp, err := m.client.Post(
		fmt.Sprintf("http://%s:80/post", ip),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("ошибка при отправке данных: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка сервера: %d", resp.StatusCode)
	}

	return nil
}
