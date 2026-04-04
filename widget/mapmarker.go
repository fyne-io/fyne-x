package widget

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MapMarker interface {
	Lat() float64
	Lon() float64
	Title() string
}

type genericMapMarker struct {
	lat   float64
	lon   float64
	title string
}

func NewMapMarker(lat, lon float64, title string) MapMarker {
	return &genericMapMarker{lat, lon, title}
}

func (m *genericMapMarker) Lat() float64 {
	return m.lat
}

func (m *genericMapMarker) Lon() float64 {
	return m.lon
}

func (m *genericMapMarker) Title() string {
	return m.title
}

type mapMarker struct {
	widget.BaseWidget
	container *fyne.Container
	title     *canvas.Text
	obj       MapMarker
	item      *mapMarkerItem
}

//go:embed mapmarker.svg
var resourceMapmarkerSvgData []byte

var resourceMapmarkerSvg = &fyne.StaticResource{
	StaticName:    "mapmarker.svg",
	StaticContent: resourceMapmarkerSvgData,
}

func newMapMarker(obj MapMarker) *mapMarker {
	m := &mapMarker{
		obj: obj,
	}
	m.ExtendBaseWidget(m)
	return m
}

func (m *mapMarker) setup() {
	if m.container != nil {
		return
	}

	m.title = canvas.NewText(m.obj.Title(), theme.ColorForWidget(theme.ColorNamePrimary, m))
	m.title.Hide()

	m.item = newMapMarkerItem(func() {
		if m.title.Hidden {
			m.title.Show()
		} else {
			m.title.Hide()
		}
		m.container.Refresh()
	})

	m.container = container.NewVBox(
		m.title,
		container.NewHBox(layout.NewSpacer(), m.item, layout.NewSpacer()),
	)
}

func (m *mapMarker) pinOffset() fyne.Position {
	m.setup()
	imgSize := m.item.Size()
	if m.title.Hidden {
		return fyne.NewPos(imgSize.Width/2, imgSize.Height+(imgSize.Height*.05))
	}
	conSize := m.container.Size()
	return fyne.NewPos(conSize.Width/2, conSize.Height+(imgSize.Height*.05))
}

func (m *mapMarker) MinSize() fyne.Size {
	m.setup()
	return m.container.MinSize()
}

func (m *mapMarker) CreateRenderer() fyne.WidgetRenderer {
	m.setup()
	return widget.NewSimpleRenderer(m.container)
}
