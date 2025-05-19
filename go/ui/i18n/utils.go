package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

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

// GetSystemLanguage возвращает язык системы
func GetSystemLanguage() string {
	// Получаем тег языка системы
	tag, _ := language.Parse(language.MustParse("en").String())

	// Если язык русский, возвращаем "ru"
	if tag == language.Russian {
		return "ru"
	}

	// По умолчанию английский
	return "en"
}
