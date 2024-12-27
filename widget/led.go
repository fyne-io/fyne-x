package widget

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var defaultLedOnColor color.Color = color.RGBA{255, 0, 0, 255}
var defaultLedOffColor color.Color = color.RGBA{120, 120, 120, 255}

type ledRenderer struct {
	led    *Led
	circle *canvas.Circle
	text   *canvas.Text
}

func (h *ledRenderer) MinSize() fyne.Size {
	log.Println(h.led.Size())
	return fyne.NewSize(h.led.Size().Width, h.led.Size().Height)
}

func (h *ledRenderer) Layout(_ fyne.Size) {
	h.circle.Resize(fyne.NewSize(h.led.MinSize().Width, h.led.MinSize().Width))
	pos := fyne.NewPos(
		(h.led.MinSize().Width-h.text.MinSize().Width)/2,
		h.led.MinSize().Height-h.text.MinSize().Height,
	)
	h.text.Move(pos)
}

func (h *ledRenderer) Destroy() {
}

func (h *ledRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{h.circle, h.text}
}

func (h *ledRenderer) Refresh() {
	h.circle.FillColor = h.led.getLedColor()
	canvas.Refresh(h.circle)
}

type Led struct {
	widget.BaseWidget
	state    bool
	text     string
	OnColor  color.Color
	offColor color.Color
}

// NewMap creates a new instance of the map widget.
func NewLed(text string) *Led {
	m := &Led{
		state:    false,
		text:     text,
		OnColor:  defaultLedOnColor,
		offColor: defaultLedOffColor,
	}
	m.ExtendBaseWidget(m)
	return m
}

func (w *Led) getLedColor() color.Color {
	if w.state {
		return w.OnColor
	} else {
		return w.offColor
	}
}

func (w *Led) MinSize() fyne.Size {
	return fyne.NewSize(30, 45)
}

func (w *Led) Set(state bool) {
	w.state = state
	w.Refresh()
}

// CreateRenderer returns the renderer for this widget.
// A map renderer is simply the map Raster with user interface elements overlaid.
func (m *Led) CreateRenderer() fyne.WidgetRenderer {
	circ0 := canvas.NewCircle(defaultLedOffColor)
	circ0.StrokeColor = color.Black
	circ0.StrokeWidth = 3

	text0 := canvas.NewText(m.text, color.Black)
	text0.TextSize = 10

	r := &ledRenderer{led: m, circle: circ0, text: text0}
	return r
}

type LedBar struct {
	widget.BaseWidget
	texts []string
	leds  []fyne.CanvasObject
}

func NewLedBar(texts []string) *LedBar {
	m := &LedBar{texts: texts}
	m.ExtendBaseWidget(m)
	m.makeLeds()
	return m
}

func (m *LedBar) makeLeds() {
	m.leds = make([]fyne.CanvasObject, len(m.texts))
	for i := 0; i < len(m.texts); i++ {
		m.leds[i] = NewLed(m.texts[i])
	}
}

func (m *LedBar) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewHBox(m.leds...)
	return widget.NewSimpleRenderer(c)
}

func (m *LedBar) SetOnColor(c color.Color) {
	for i := 0; i < len(m.leds); i++ {
		m.leds[i].(*Led).OnColor = c
	}
}

func (w *LedBar) Set(state int) {

	var l int = len(w.texts) - 1
	var b int = 1

	for i := 0; i < len(w.texts); i++ {
		if (state & b) == b {
			w.leds[l-i].(*Led).Set(true)
		} else {
			w.leds[l-i].(*Led).Set(false)
		}
		b <<= 1
	}

	w.Refresh()
}
