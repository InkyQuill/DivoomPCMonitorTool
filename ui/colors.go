package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CardColor возвращает цвет фона для карточки
func CardColor() color.Color {
	if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantDark {
		return color.NRGBA{40, 40, 40, 255}
	}
	return color.NRGBA{240, 240, 240, 255}
}

// MiniBoxColor возвращает цвет фона для мини-бокса
func MiniBoxColor() color.Color {
	if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantDark {
		return color.NRGBA{60, 60, 60, 255}
	}
	return color.NRGBA{230, 230, 230, 255}
}
