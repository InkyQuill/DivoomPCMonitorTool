package tray

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

// InitTray инициализирует системный трей
func InitTray(a fyne.App, w fyne.Window) {
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("System Metrics",
			fyne.NewMenuItem("Show", func() {
				w.Show()
			}),
			fyne.NewMenuItem("Hide", func() {
				w.Hide()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() {
				a.Quit()
			}),
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

		// Скрываем окно при закрытии или сворачивании
		w.SetOnClosed(func() {
			w.Hide()
		})
	}
}
