package widget

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

func NewMapMarker(lat, lon float64, title string) *genericMapMarker {
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
	obj MapMarker
	img *canvas.Image
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
	m.img = canvas.NewImageFromResource(theme.NewColoredResource(resourceMapmarkerSvg, theme.ColorNamePrimary))
	m.img.SetMinSize(fyne.NewSize(50, 50))
}

func (m *mapMarker) MinSize() fyne.Size {
	if m.img == nil {
		m.setup()
	}
	return m.img.MinSize()
}

func (m *mapMarker) CreateRenderer() fyne.WidgetRenderer {
	if m.img == nil {
		m.setup()
	}
	return widget.NewSimpleRenderer(m.img)
}
