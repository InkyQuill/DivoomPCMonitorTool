package ui

import (
	"divoom/ui/i18n"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// InfoPage содержит все компоненты страницы информации
type InfoPage struct {
	Content      fyne.CanvasObject
	MetricLabels *MetricLabels
	Translator   *i18n.Translator
}

// CreateInfoPage создает страницу с информацией о системе
func CreateInfoPage(sysEdition, sysKernel, cpuModel, gpuModel string) *InfoPage {
	translator := i18n.NewTranslator("en") // Временный переводчик, будет заменен в main.go
	metricLabels := NewMetricLabels()

	// Метки для CPU
	cpuModelLabel := NewCustomLabel(cpuModel, 16, true, false)
	cpuTempUnit := NewUnitLabel("°C")
	cpuUsageUnit := NewUnitLabel("%")

	// Метки для GPU
	gpuModelLabel := NewCustomLabel(gpuModel, 16, true, false)
	gpuTempUnit := NewUnitLabel("°C")
	gpuUsageUnit := NewUnitLabel("%")

	// RAM
	ramUnit := NewUnitLabel("%")

	// HDD
	hddTempUnit := NewUnitLabel("°C")

	// --- Блоки метрик ---
	cpuBox := NewMetricCard(translator.T("metric.cpu"), cpuModelLabel,
		NewMiniMetricBox(translator.T("metric.temp"), metricLabels.CPUTemp, cpuTempUnit),
		NewMiniMetricBox(translator.T("metric.usage"), metricLabels.CPUUsage, cpuUsageUnit),
	)
	gpuBox := NewMetricCard(translator.T("metric.gpu"), gpuModelLabel,
		NewMiniMetricBox(translator.T("metric.temp"), metricLabels.GPUTemp, gpuTempUnit),
		NewMiniMetricBox(translator.T("metric.usage"), metricLabels.GPUUsage, gpuUsageUnit),
	)
	ramBox := NewSimpleMetricCard(translator.T("metric.ram"), metricLabels.RAMUsage, ramUnit)
	hddBox := NewSimpleMetricCard(translator.T("metric.hdd"), metricLabels.HDDTemp, hddTempUnit)

	// --- Верхняя панель ---
	sysEditionLabel := NewCustomLabel(sysEdition, 16, false, false)
	sysKernelLabel := NewCustomLabel(translator.T("system.kernel", sysKernel), 12, false, false)

	timeGateTitle := NewCustomLabel(translator.T("system.timegate"), 16, true, true)
	timeGateIP := NewCustomLabel("192.168.2.21", 12, false, true)

	topPanel := container.New(NewPaddedLayout(12), container.NewHBox(
		container.NewVBox(
			sysEditionLabel,
			sysKernelLabel,
		),
		layout.NewSpacer(),
		container.NewVBox(
			timeGateTitle,
			timeGateIP,
		),
	))

	// --- Основной layout ---
	content := container.New(NewPaddedLayout(12), container.NewVBox(
		topPanel,
		container.NewGridWithColumns(2, cpuBox, gpuBox),
		container.NewGridWithColumns(2, ramBox, hddBox),
	))

	return &InfoPage{
		Content:      content,
		MetricLabels: metricLabels,
		Translator:   translator,
	}
}

// SettingsPage содержит все компоненты страницы настроек
type SettingsPage struct {
	Content        fyne.CanvasObject
	AutostartCheck *widget.Check
	MinimizedCheck *widget.Check
	LanguageSelect *widget.Select
	Translator     *i18n.Translator
}

// CreateSettingsPage создает страницу настроек
func CreateSettingsPage(translator *i18n.Translator) *SettingsPage {
	autostartCheck := widget.NewCheck(translator.T("settings.autostart"), nil)
	minimizedCheck := widget.NewCheck(translator.T("settings.minimized"), nil)

	languageSelect := widget.NewSelect([]string{"English", "Русский"}, nil)
	// Устанавливаем выбранный язык в зависимости от текущего переводчика
	if translator.Lang == "ru" {
		languageSelect.SetSelected("Русский")
	} else {
		languageSelect.SetSelected("English")
	}

	content := container.New(NewPaddedLayout(12),
		container.NewVBox(
			widget.NewLabel(translator.T("settings.title")),
			autostartCheck,
			minimizedCheck,
			widget.NewLabel(translator.T("settings.language")),
			languageSelect,
		),
	)

	return &SettingsPage{
		Content:        content,
		AutostartCheck: autostartCheck,
		MinimizedCheck: minimizedCheck,
		LanguageSelect: languageSelect,
		Translator:     translator,
	}
}
