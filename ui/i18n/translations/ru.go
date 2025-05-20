package translations

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// RegisterRussian регистрирует русские переводы
func RegisterRussian() {
	// App
	message.SetString(language.Russian, "app.title", "Divoom PC Monitor")

	// Tabs
	message.SetString(language.Russian, "tab.info", "Информация")
	message.SetString(language.Russian, "tab.settings", "Настройки")

	// Metrics
	message.SetString(language.Russian, "metric.cpu", "Процессор")
	message.SetString(language.Russian, "metric.gpu", "Видеокарта")
	message.SetString(language.Russian, "metric.ram", "Память")
	message.SetString(language.Russian, "metric.hdd", "Диск")
	message.SetString(language.Russian, "metric.temp", "Температура")
	message.SetString(language.Russian, "metric.usage", "Загрузка")

	// System
	message.SetString(language.Russian, "system.kernel", "Ядро: %s")
	message.SetString(language.Russian, "system.device", "Устройство")

	// Settings
	message.SetString(language.Russian, "settings.title", "Настройки")
	message.SetString(language.Russian, "settings.autostart", "Запускать при старте системы")
	message.SetString(language.Russian, "settings.minimized", "Запускать свёрнутым")
	message.SetString(language.Russian, "settings.language", "Язык")
	message.SetString(language.Russian, "settings.device", "Устройство")
	message.SetString(language.Russian, "settings.no_devices", "Устройства не найдены")

	// Tray Menu
	message.SetString(language.Russian, "tray.show", "Показать окно")
	message.SetString(language.Russian, "tray.hide", "Спрятать окно")
	message.SetString(language.Russian, "tray.quit", "Выйти")
}
