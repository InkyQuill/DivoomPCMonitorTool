package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"divoom/autostart"
)

const (
	configDir  = ".divoom"
	configFile = "config.json"
)

// DeviceInfo содержит информацию о выбранном устройстве
type DeviceInfo struct {
	IP          string `json:"ip"`
	Type        int    `json:"type"`
	SelectedLCD int    `json:"selected_lcd"`
}

// Settings содержит все настройки приложения
type Settings struct {
	StartMinimized bool       `json:"start_minimized"`
	Language       string     `json:"language"` // "en" или "ru"
	Device         DeviceInfo `json:"device"`
	// Здесь будут добавляться новые настройки
}

// Manager управляет настройками приложения
type Manager struct {
	settings *Settings
	path     string
}

// NewManager создает новый менеджер настроек
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	configPath := filepath.Join(homeDir, configDir)
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	manager := &Manager{
		settings: &Settings{
			StartMinimized: false,
			Language:       "en", // По умолчанию английский
		},
		path: filepath.Join(configPath, configFile),
	}

	// Загружаем существующие настройки
	if err := manager.Load(); err != nil {
		return nil, fmt.Errorf("failed to load settings: %v", err)
	}

	return manager, nil
}

// Load загружает настройки из файла
func (m *Manager) Load() error {
	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			// Если файл не существует, создаем его с настройками по умолчанию
			return m.Save()
		}
		return fmt.Errorf("failed to read config file: %v", err)
	}

	if err := json.Unmarshal(data, m.settings); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	return nil
}

// Save сохраняет настройки в файл
func (m *Manager) Save() error {
	data, err := json.MarshalIndent(m.settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %v", err)
	}

	if err := os.WriteFile(m.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// Get возвращает текущие настройки
func (m *Manager) Get() *Settings {
	return m.settings
}

// SetStartMinimized устанавливает флаг запуска свёрнутым
func (m *Manager) SetStartMinimized(value bool) error {
	m.settings.StartMinimized = value
	return m.Save()
}

// SetLanguage устанавливает язык интерфейса
func (m *Manager) SetLanguage(lang string) error {
	if lang != "en" && lang != "ru" {
		return fmt.Errorf("unsupported language: %s", lang)
	}
	m.settings.Language = lang
	return m.Save()
}

// SetDevice устанавливает информацию об устройстве
func (m *Manager) SetDevice(ip string, deviceType int, selectedLCD int) error {
	m.settings.Device = DeviceInfo{
		IP:          ip,
		Type:        deviceType,
		SelectedLCD: selectedLCD,
	}
	return m.Save()
}

// GetDevice возвращает информацию о выбранном устройстве
func (m *Manager) GetDevice() DeviceInfo {
	return m.settings.Device
}

// IsAutostartEnabled проверяет, включен ли автозапуск
func (m *Manager) IsAutostartEnabled() bool {
	return autostart.IsEnabled()
}

// SetAutostart устанавливает флаг автозапуска
func (m *Manager) SetAutostart(value bool) error {
	if value {
		return autostart.Enable()
	}
	return autostart.Disable()
}
