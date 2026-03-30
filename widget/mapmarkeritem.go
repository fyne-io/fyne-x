package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type mapMarkerItem struct {
	widget.BaseWidget
	img *canvas.Image
	fn  func()
}

func newMapMarkerItem(fn func()) *mapMarkerItem {
	item := &mapMarkerItem{fn: fn}
	item.ExtendBaseWidget(item)
	return item
}

func (m *mapMarkerItem) Tapped(ev *fyne.PointEvent) {
	if m.fn == nil {
		return
	}
	m.fn()
}

func (m *mapMarkerItem) setup() {
	if m.img != nil {
		return
	}
	res := theme.NewColoredResource(resourceMapmarkerSvg, theme.ColorNamePrimary)
	m.img = canvas.NewImageFromResource(res)
	m.img.SetMinSize(fyne.NewSize(50, 50))
}

func (m *mapMarkerItem) CreateRenderer() fyne.WidgetRenderer {
	m.setup()
	return widget.NewSimpleRenderer(m.img)
}
