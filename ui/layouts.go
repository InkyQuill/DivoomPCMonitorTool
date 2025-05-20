package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// NewPaddedLayout создает layout с отступами
func NewPaddedLayout(padding float32) fyne.Layout {
	return layout.NewCustomPaddedLayout(padding, padding, padding, padding)
}

// NewMetricCard создает карточку с метриками
func NewMetricCard(title string, model fyne.CanvasObject, left, right fyne.CanvasObject) fyne.CanvasObject {
	bg := NewCardBackground(220, 110)
	return container.NewStack(
		bg,
		container.New(NewPaddedLayout(12),
			container.NewVBox(
				NewCustomLabel(title, 12, true, false),
				model,
				container.NewHBox(
					left,
					layout.NewSpacer(),
					right,
				),
			)),
	)
}

// NewSimpleMetricCard создает простую карточку с одной метрикой
func NewSimpleMetricCard(title string, value, unit fyne.CanvasObject) fyne.CanvasObject {
	bg := NewCardBackground(220, 60)
	return container.NewStack(
		bg,
		container.New(NewPaddedLayout(12),
			container.NewVBox(
				NewCustomLabel(title, 12, true, false),
				container.NewHBox(
					layout.NewSpacer(),
					container.NewHBox(value, unit),
				),
			)),
	)
}

// NewMiniMetricBox создает мини-бокс для отображения метрики
func NewMiniMetricBox(title string, value, unit fyne.CanvasObject) fyne.CanvasObject {
	bg := NewMiniBoxBackground()
	return container.NewStack(
		bg,
		container.New(NewPaddedLayout(8),
			container.NewVBox(
				NewCustomLabel(title, 10, false, false),
				container.NewHBox(layout.NewSpacer(), value, unit),
			)),
	)
}
