package i18n

import (
	"os"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"divoom/ui/i18n/translations"
)

var (
	// Инициализируем принтеры для разных языков
	enPrinter *message.Printer
	ruPrinter *message.Printer
)

// Инициализируем переводы
func init() {
	// Регистрируем переводы
	translations.RegisterEnglish()
	translations.RegisterRussian()

	// Создаем принтеры
	enPrinter = message.NewPrinter(language.English)
	ruPrinter = message.NewPrinter(language.Russian)
}

// Translator представляет собой переводчик
type Translator struct {
	printer *message.Printer
	Lang    string
}

// NewTranslator создает новый переводчик для указанного языка
func NewTranslator(lang string) *Translator {
	return &Translator{
		printer: GetPrinter(lang),
		Lang:    lang,
	}
}

// T возвращает локализованную строку
func (t *Translator) T(key string, args ...interface{}) string {
	return t.printer.Sprintf(key, args...)
}

// GetPrinter возвращает принтер для указанного языка
func GetPrinter(lang string) *message.Printer {
	if lang == "ru" {
		return ruPrinter
	}
	return enPrinter
}

// GetSystemLanguage возвращает язык системы
func GetSystemLanguage() string {
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LANGUAGE")
	}
	if lang == "" {
		return "en"
	}

	// Извлекаем код языка (например, "ru_RU.UTF-8" -> "ru")
	lang = strings.Split(lang, "_")[0]
	if lang == "ru" {
		return "ru"
	}
	return "en"
}
