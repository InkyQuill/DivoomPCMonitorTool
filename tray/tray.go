package tray

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"

	"divoom/ui/i18n"
)

// InitTray инициализирует системный трей
func InitTray(a fyne.App, w fyne.Window) {
	if desk, ok := a.(desktop.App); ok {
		translator := i18n.NewTranslator(i18n.GetSystemLanguage())

		// Создаем меню
		showItem := fyne.NewMenuItem(translator.T("tray.show"), func() {
			w.Show()
		})
		hideItem := fyne.NewMenuItem(translator.T("tray.hide"), func() {
			w.Hide()
		})
		quitItem := fyne.NewMenuItem(translator.T("tray.quit"), func() {
			a.Quit()
		})

		// Обновляем меню при изменении видимости окна
		w.SetOnClosed(func() {
			w.Hide()
		})

		m := fyne.NewMenu("System Metrics",
			showItem,
			hideItem,
			quitItem,
		)
		desk.SetSystemTrayMenu(m)

		// Загружаем иконку для трея
		icon, err := fyne.LoadResourceFromPath("heart.png")
		if err != nil {
			// Если не удалось загрузить иконку, используем стандартную
			desk.SetSystemTrayIcon(theme.ComputerIcon())
		} else {
			desk.SetSystemTrayIcon(icon)
		}
	}
}
