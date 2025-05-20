package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

// getTextColor возвращает цвет текста в зависимости от темы
func getTextColor() color.Color {
	if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantDark {
		return color.White
	}
	return color.Black
}

// NewCustomLabel создает текстовую метку с заданными параметрами
func NewCustomLabel(text string, size int, bold, right bool) *canvas.Text {
	label := canvas.NewText(text, getTextColor())
	label.TextSize = float32(size)
	label.TextStyle.Bold = bold
	if right {
		label.Alignment = fyne.TextAlignTrailing
	}
	return label
}

// NewValueLabel создает метку для отображения значения
func NewValueLabel(size int) *canvas.Text {
	label := canvas.NewText("0", getTextColor())
	label.TextSize = float32(size)
	label.TextStyle.Bold = true
	return label
}

// NewUnitLabel создает метку для отображения единиц измерения
func NewUnitLabel(unit string) *canvas.Text {
	label := canvas.NewText(unit, getTextColor())
	label.TextSize = 12
	return label
}

// NewCardBackground создает фон для карточки
func NewCardBackground(width, height float32) *canvas.Rectangle {
	bg := canvas.NewRectangle(CardColor())
	bg.SetMinSize(fyne.NewSize(width, height))
	bg.StrokeColor = color.NRGBA{0, 0, 0, 32}
	bg.StrokeWidth = 1
	bg.CornerRadius = 12
	return bg
}

// NewMiniBoxBackground создает фон для мини-бокса
func NewMiniBoxBackground() *canvas.Rectangle {
	bg := canvas.NewRectangle(MiniBoxColor())
	bg.SetMinSize(fyne.NewSize(160, 40))
	bg.StrokeColor = color.NRGBA{0, 0, 0, 16}
	bg.StrokeWidth = 1
	bg.CornerRadius = 8
	return bg
}
