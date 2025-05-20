package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"divoom/device"
	"divoom/logger"
	"divoom/metrics"
	"divoom/settings"
	"divoom/tray"
	"divoom/ui"
	"divoom/ui/i18n"
)

// copyResource копирует ресурс в директорию приложения
func copyResource(src, dst string) error {
	// Получаем путь к исполняемому файлу
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exe)

	// Создаем полный путь назначения
	dstPath := filepath.Join(exeDir, dst)

	// Копируем файл
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dstPath, data, 0644)
}

func main() {
	// Инициализируем логгер
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
	}
	defer logger.Close()

	// Инициализируем менеджер настроек
	settingsManager, err := settings.NewManager()
	if err != nil {
		fmt.Printf("Failed to initialize settings: %v\n", err)
		return
	}

	// Инициализируем менеджер устройств
	deviceManager := device.NewManager()

	// Определяем язык интерфейса
	lang := settingsManager.Get().Language
	if lang == "" {
		// Если язык не установлен, используем системный
		lang = i18n.GetSystemLanguage()
		settingsManager.SetLanguage(lang)
	}

	// Создаем приложение
	myApp := app.NewWithID("com.divoom.pcmonitor")
	myWindow := myApp.NewWindow("Divoom PC Monitor") // Временный заголовок

	// Канал для сигнала о смене языка
	langChangeChan := make(chan string)

	// Функция для создания интерфейса
	createUI := func(lang string) {
		// Создаем переводчик
		translator := i18n.NewTranslator(lang)

		// Обновляем заголовок окна
		myWindow.SetTitle(translator.T("app.title"))

		// Получаем информацию о системе
		var sysEdition, sysKernel string
		switch runtime.GOOS {
		case "linux":
			distro, kernel := metrics.GetLinuxSystemInfo()
			sysEdition = distro
			sysKernel = kernel
		case "windows":
			edition, build, arch := metrics.GetWindowsSystemInfo()
			sysEdition = edition + " (" + arch + ")"
			sysKernel = build
		case "darwin":
			edition := metrics.GetDarwinSystemInfo()
			sysEdition = edition
			sysKernel = ""
		default:
			sysEdition = "Unknown system"
			sysKernel = ""
		}

		// Создаем страницы
		infoPage := ui.CreateInfoPage(sysEdition, sysKernel, getCPUModel(), getGPUModel(), settingsManager.GetDevice())
		infoPage.Translator = translator // Обновляем переводчик
		settingsPage := ui.CreateSettingsPage(translator, deviceManager, settingsManager)

		// Настраиваем обработчики событий для настроек
		settingsPage.AutostartCheck.SetChecked(settingsManager.IsAutostartEnabled())
		settingsPage.AutostartCheck.OnChanged = func(checked bool) {
			settingsManager.SetAutostart(checked)
		}
		settingsPage.MinimizedCheck.SetChecked(settingsManager.Get().StartMinimized)
		settingsPage.MinimizedCheck.OnChanged = func(checked bool) {
			settingsManager.SetStartMinimized(checked)
		}
		settingsPage.LanguageSelect.OnChanged = func(selected string) {
			var newLang string
			switch selected {
			case "Русский":
				newLang = "ru"
			default:
				newLang = "en"
			}
			settingsManager.SetLanguage(newLang)
			// Отправляем сигнал о смене языка
			langChangeChan <- newLang
		}

		// Создаем вкладки
		tabs := container.NewAppTabs(
			container.NewTabItem(translator.T("tab.info"), infoPage.Content),
			container.NewTabItem(translator.T("tab.settings"), settingsPage.Content),
		)
		tabs.SetTabLocation(container.TabLocationTop)

		myWindow.SetContent(tabs)

		// Запускаем сбор метрик и отправку на устройство
		collector := metrics.NewCollector()
		go func() {
			for {
				if err := collector.Update(); err != nil {
					fmt.Printf("Error updating metrics: %v\n", err)
					continue
				}
				metrics := collector.GetMetrics()

				fyne.Do(func() {
					infoPage.MetricLabels.UpdateMetrics(
						metrics.CPUTemp,
						metrics.CPUUsage,
						metrics.GPUTemp,
						metrics.GPUUsage,
						metrics.RAMUsage,
						metrics.HDDTemp,
					)
				})

				// Отправляем метрики на устройство
				deviceInfo := settingsManager.GetDevice()
				if deviceInfo.IP != "" {
					metricsMap := map[string]float64{
						"cpu_usage":    metrics.CPUUsage,
						"gpu_usage":    metrics.GPUUsage,
						"cpu_temp":     metrics.CPUTemp,
						"gpu_temp":     metrics.GPUTemp,
						"memory_usage": metrics.RAMUsage,
						"disk_usage":   metrics.HDDTemp,
					}
					if err := deviceManager.SendSystemInfo(deviceInfo.IP, deviceInfo.SelectedLCD, metricsMap); err != nil {
						logger.Log("Error sending metrics to device: %v", err)
					}
				}

				time.Sleep(time.Second)
			}
		}()
	}

	// Создаем начальный интерфейс
	createUI(lang)

	// Запускаем горутину для обработки смены языка
	go func() {
		for newLang := range langChangeChan {
			fyne.Do(func() {
				createUI(newLang)
			})
		}
	}()

	// Копируем ресурсы
	if err := copyResource("Divoom.ico", "Divoom.ico"); err != nil {
		logger.Log("Error copying icon: %v", err)
	}
	if err := copyResource("heart.png", "heart.png"); err != nil {
		logger.Log("Error copying heart icon: %v", err)
	}

	// Загружаем иконку
	icon, err := fyne.LoadResourceFromPath("heart.png")
	if err != nil {
		logger.Log("Error loading icon: %v", err)
	} else {
		myWindow.SetIcon(icon)
		myApp.SetIcon(icon)
	}

	myWindow.Resize(fyne.NewSize(500, 500))
	tray.InitTray(myApp, myWindow)

	// Если установлен флаг запуска свёрнутым, скрываем окно
	if settingsManager.Get().StartMinimized {
		myWindow.Hide()
	} else {
		myWindow.Show()
	}

	myApp.Run()
}

// getCPUModel возвращает модель CPU
func getCPUModel() string {
	return metrics.GetCPUInfo()
}

// getGPUModel возвращает модель GPU
func getGPUModel() string {
	return metrics.GetGPUInfo()
}
