package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Инициализируем принтеры для разных языков
var (
	enPrinter *message.Printer
	ruPrinter *message.Printer
)

// Инициализируем сообщения
func init() {
	enPrinter = message.NewPrinter(language.English)
	ruPrinter = message.NewPrinter(language.Russian)

	// Заголовок приложения
	message.SetString(language.English, "app.title", "Divoom PC Monitor")
	message.SetString(language.Russian, "app.title", "Divoom PC Monitor")

	// Вкладки
	message.SetString(language.English, "tab.info", "Info")
	message.SetString(language.Russian, "tab.info", "Информация")
	message.SetString(language.English, "tab.settings", "Settings")
	message.SetString(language.Russian, "tab.settings", "Настройки")

	// Метрики
	message.SetString(language.English, "metric.cpu", "CPU")
	message.SetString(language.Russian, "metric.cpu", "Процессор")
	message.SetString(language.English, "metric.gpu", "GPU")
	message.SetString(language.Russian, "metric.gpu", "Видеокарта")
	message.SetString(language.English, "metric.ram", "RAM")
	message.SetString(language.Russian, "metric.ram", "Память")
	message.SetString(language.English, "metric.hdd", "HDD")
	message.SetString(language.Russian, "metric.hdd", "Диск")

	// Подметрики
	message.SetString(language.English, "metric.temp", "Temperature")
	message.SetString(language.Russian, "metric.temp", "Температура")
	message.SetString(language.English, "metric.usage", "Usage")
	message.SetString(language.Russian, "metric.usage", "Использование")

	// Настройки
	message.SetString(language.English, "settings.title", "Settings")
	message.SetString(language.Russian, "settings.title", "Настройки")
	message.SetString(language.English, "settings.autostart", "Autostart on system boot")
	message.SetString(language.Russian, "settings.autostart", "Автозапуск при старте системы")
	message.SetString(language.English, "settings.minimized", "Start minimized")
	message.SetString(language.Russian, "settings.minimized", "Запуск свёрнутым")
	message.SetString(language.English, "settings.language", "Interface language")
	message.SetString(language.Russian, "settings.language", "Язык интерфейса")

	// Системная информация
	message.SetString(language.English, "system.kernel", "Kernel: %s")
	message.SetString(language.Russian, "system.kernel", "Ядро: %s")
	message.SetString(language.English, "system.timegate", "Time Gate")
	message.SetString(language.Russian, "system.timegate", "Time Gate")
}

// GetPrinter возвращает принтер для указанного языка
func GetPrinter(lang string) *message.Printer {
	if lang == "ru" {
		return ruPrinter
	}
	return enPrinter
}
