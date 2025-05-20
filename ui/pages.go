package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"divoom/device"
	"divoom/settings"
	"divoom/ui/i18n"
)

// InfoPage содержит все компоненты страницы информации
type InfoPage struct {
	Content      fyne.CanvasObject
	MetricLabels *MetricLabels
	Translator   *i18n.Translator
	DeviceLabel  *widget.Label
}

// CreateInfoPage создает страницу с информацией о системе
func CreateInfoPage(sysEdition, sysKernel, cpuModel, gpuModel string, deviceInfo settings.DeviceInfo) *InfoPage {
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

	deviceLabel := widget.NewLabel("")
	if deviceInfo.IP != "" {
		deviceLabel.SetText(fmt.Sprintf("%s: %s", translator.T("system.device"), deviceInfo.IP))
	}

	topPanel := container.New(NewPaddedLayout(12), container.NewHBox(
		container.NewVBox(
			sysEditionLabel,
			sysKernelLabel,
		),
		layout.NewSpacer(),
		container.NewVBox(
			deviceLabel,
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
		DeviceLabel:  deviceLabel,
	}
}

// SettingsPage содержит все компоненты страницы настроек
type SettingsPage struct {
	Content         fyne.CanvasObject
	AutostartCheck  *widget.Check
	MinimizedCheck  *widget.Check
	LanguageSelect  *widget.Select
	DeviceSelect    *widget.Select
	LCDSelect       *widget.Select
	Translator      *i18n.Translator
	DeviceManager   *device.Manager
	SettingsManager *settings.Manager
}

// CreateSettingsPage создает страницу настроек
func CreateSettingsPage(translator *i18n.Translator, deviceManager *device.Manager, settingsManager *settings.Manager) *SettingsPage {
	autostartCheck := widget.NewCheck(translator.T("settings.autostart"), nil)
	minimizedCheck := widget.NewCheck(translator.T("settings.minimized"), nil)

	languageSelect := widget.NewSelect([]string{"English", "Русский"}, nil)
	if translator.Lang == "ru" {
		languageSelect.SetSelected("Русский")
	} else {
		languageSelect.SetSelected("English")
	}

	deviceSelect := widget.NewSelect([]string{}, nil)
	lcdSelect := widget.NewSelect([]string{}, nil)
	lcdSelect.Hide()

	// Функция для обновления списка устройств
	updateDevices := func() {
		devices, err := deviceManager.DiscoverDevices()
		if err != nil {
			deviceSelect.Options = []string{translator.T("settings.no_devices")}
			deviceSelect.SetSelected(translator.T("settings.no_devices"))
			return
		}

		if len(devices) == 0 {
			deviceSelect.Options = []string{translator.T("settings.no_devices")}
			deviceSelect.SetSelected(translator.T("settings.no_devices"))
			return
		}

		options := make([]string, len(devices))
		for i, d := range devices {
			options[i] = fmt.Sprintf("%s (%s)", d.DeviceName, d.DevicePrivateIP)
		}
		deviceSelect.Options = options

		// Выбираем сохраненное устройство
		savedDevice := settingsManager.GetDevice()
		if savedDevice.IP != "" {
			for i, d := range devices {
				if d.DevicePrivateIP == savedDevice.IP {
					deviceSelect.SetSelected(options[i])
					if d.Hardware == device.DeviceTypeTimeGate {
						lcdSelect.Show()
						lcdSelect.SetSelected(fmt.Sprintf("%d", savedDevice.SelectedLCD+1))
					}
					break
				}
			}
		}
	}

	// Обработчик выбора устройства
	deviceSelect.OnChanged = func(selected string) {
		if selected == translator.T("settings.no_devices") {
			lcdSelect.Hide()
			return
		}

		devices, _ := deviceManager.DiscoverDevices()
		for _, d := range devices {
			if fmt.Sprintf("%s (%s)", d.DeviceName, d.DevicePrivateIP) == selected {
				if d.Hardware == device.DeviceTypeTimeGate {
					lcdSelect.Options = []string{"1", "2", "3", "4", "5"}
					lcdSelect.Show()
					lcdSelect.SetSelected("1")
				} else {
					lcdSelect.Hide()
				}
				settingsManager.SetDevice(d.DevicePrivateIP, d.Hardware, 0)
				break
			}
		}
	}

	// Обработчик выбора экрана
	lcdSelect.OnChanged = func(selected string) {
		if selected == "" {
			return
		}
		lcdID := 0
		fmt.Sscanf(selected, "%d", &lcdID)
		devices, _ := deviceManager.DiscoverDevices()
		for _, d := range devices {
			if d.Hardware == device.DeviceTypeTimeGate {
				settingsManager.SetDevice(d.DevicePrivateIP, d.Hardware, lcdID-1)
				break
			}
		}
	}

	// Обновляем список устройств
	updateDevices()

	content := container.New(NewPaddedLayout(12),
		container.NewVBox(
			widget.NewLabel(translator.T("settings.title")),
			autostartCheck,
			minimizedCheck,
			widget.NewLabel(translator.T("settings.language")),
			languageSelect,
			widget.NewLabel(translator.T("settings.device")),
			deviceSelect,
			lcdSelect,
		),
	)

	return &SettingsPage{
		Content:         content,
		AutostartCheck:  autostartCheck,
		MinimizedCheck:  minimizedCheck,
		LanguageSelect:  languageSelect,
		DeviceSelect:    deviceSelect,
		LCDSelect:       lcdSelect,
		Translator:      translator,
		DeviceManager:   deviceManager,
		SettingsManager: settingsManager,
	}
}
