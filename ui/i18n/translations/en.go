package translations

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// RegisterEnglish регистрирует английские переводы
func RegisterEnglish() {
	// App
	message.SetString(language.English, "app.title", "Divoom PC Monitor")

	// Tabs
	message.SetString(language.English, "tab.info", "System Info")
	message.SetString(language.English, "tab.settings", "Settings")

	// Metrics
	message.SetString(language.English, "metric.cpu", "CPU")
	message.SetString(language.English, "metric.gpu", "GPU")
	message.SetString(language.English, "metric.ram", "RAM")
	message.SetString(language.English, "metric.hdd", "HDD")
	message.SetString(language.English, "metric.temp", "Temperature")
	message.SetString(language.English, "metric.usage", "Usage")

	// System
	message.SetString(language.English, "system.kernel", "Kernel: %s")
	message.SetString(language.English, "system.device", "Device")

	// Settings
	message.SetString(language.English, "settings.title", "Settings")
	message.SetString(language.English, "settings.autostart", "Start with system")
	message.SetString(language.English, "settings.minimized", "Start minimized")
	message.SetString(language.English, "settings.language", "Language")
	message.SetString(language.English, "settings.device", "Device")
	message.SetString(language.English, "settings.no_devices", "No devices found")

	// Tray Menu
	message.SetString(language.English, "tray.show", "Show Window")
	message.SetString(language.English, "tray.hide", "Hide Window")
	message.SetString(language.English, "tray.quit", "Quit")
}
