package widget

import (
	"image/color"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*Gauge)(nil)

const (
	gaugeSweep = math.Pi * 1.5    // the dial covers 270°
	gaugeStart = -math.Pi * 3 / 4 // needle angle at Min, radians clockwise from 12 o'clock
)

// Gauge is a widget that displays a value within a range as a needle on a
// circular dial, with graduation marks, graduation labels and a numeric readout.
type Gauge struct {
	widget.BaseWidget

	// Min and Max define the range of the gauge. They default to 0 and 100.
	Min, Max float64
	// Value is the current value indicated by the needle, use SetValue to change it.
	// The needle is clamped to the range but the readout shows the value as-is.
	Value float64
	// Title is drawn below the centre of the dial.
	Title string
	// Graduations is the number of graduation divisions around the dial, default 10.
	// Every second graduation is drawn larger and labelled with its value.
	Graduations int

	// TextFormatter can be used to have a custom format of the readout text.
	// If set, it overrides the default numeric readout and runs each time the value updates.
	TextFormatter func() string `json:"-"`

	minSize fyne.Size
}

// NewGauge creates a new gauge widget with a range of 0 to 100.
func NewGauge() *Gauge {
	g := &Gauge{Min: 0, Max: 100, minSize: fyne.NewSize(100, 100)}
	g.ExtendBaseWidget(g)
	return g
}

// SetValue changes the value indicated by this gauge (from g.Min to g.Max).
// The widget will be refreshed to indicate the change.
func (g *Gauge) SetValue(v float64) {
	g.Value = v
	g.Refresh()
}

// SetMinSize sets the size that this widget should not shrink below.
func (g *Gauge) SetMinSize(size fyne.Size) {
	g.minSize = size
	g.Refresh()
}

// MinSize returns the size that this widget should not shrink below.
func (g *Gauge) MinSize() fyne.Size {
	g.ExtendBaseWidget(g)
	return g.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (g *Gauge) CreateRenderer() fyne.WidgetRenderer {
	g.ExtendBaseWidget(g)
	if g.Min == 0 && g.Max == 0 {
		g.Max = 100
	}

	r := &gaugeRenderer{gauge: g}
	// the arc ends overshoot ±135° slightly to stay flush with the outer graduation marks
	r.face = &canvas.Arc{StartAngle: -135.73, EndAngle: 135.8, CutoutRatio: 0.985}
	r.center = &canvas.Circle{}
	r.needle = &canvas.Line{}
	r.title = &canvas.Text{Text: g.Title, Alignment: fyne.TextAlignCenter, TextStyle: fyne.TextStyle{Monospace: true}}
	r.readout = &canvas.Text{Alignment: fyne.TextAlignCenter}

	r.rebuildGraduations()
	r.applyTheme()
	r.readout.Text = r.formatValue()
	return r
}

var _ fyne.WidgetRenderer = (*gaugeRenderer)(nil)

type gaugeRenderer struct {
	gauge *Gauge

	face    *canvas.Arc
	center  *canvas.Circle
	needle  *canvas.Line
	title   *canvas.Text
	readout *canvas.Text

	graduations []*canvas.Line
	labels      []*canvas.Text // indexed as graduations, nil for minor graduations
	objects     []fyne.CanvasObject

	// the values the graduations were built from, to detect changes on Refresh
	builtMin, builtMax float64
	builtGraduations   int
	builtFg            color.Color
	graduationCount    int // effective graduation divisions, defaulted if Graduations is unset

	middle                     fyne.Position
	needleOffset, needleLength float32
}

func (r *gaugeRenderer) rebuildGraduations() {
	g := r.gauge
	r.builtMin, r.builtMax, r.builtGraduations = g.Min, g.Max, g.Graduations

	graduations := g.Graduations
	if graduations <= 0 {
		graduations = 10
	}
	r.graduationCount = graduations

	r.graduations = make([]*canvas.Line, graduations+1)
	r.labels = make([]*canvas.Text, graduations+1)
	r.objects = make([]fyne.CanvasObject, 0, 2*(graduations+1)+5)
	for i := 0; i <= graduations; i++ {
		r.graduations[i] = &canvas.Line{StrokeColor: gaugeGraduationColor(i, graduations)}
		r.objects = append(r.objects, r.graduations[i])

		if i%2 != 0 {
			continue
		}
		v := g.Min + float64(i)/float64(graduations)*(g.Max-g.Min)
		r.labels[i] = &canvas.Text{Text: strconv.FormatFloat(v, 'f', -1, 64), Alignment: fyne.TextAlignCenter}
		r.objects = append(r.objects, r.labels[i])
	}
	r.objects = append(r.objects, r.face, r.title, r.center, r.needle, r.readout)
}

// gaugeGraduationColor blends the graduation marks from green at the low end of the dial
// through yellow in the middle to red at the high end.
func gaugeGraduationColor(i, graduations int) color.Color {
	half := float64(graduations) / 2
	if float64(i) <= half {
		return color.NRGBA{R: uint8(255 * float64(i) / half), G: 0xFF, A: 0xFF}
	}
	return color.NRGBA{R: 0xFF, G: uint8(255 * (1 - (float64(i)-half)/half)), A: 0xFF}
}

func (r *gaugeRenderer) applyTheme() {
	th := r.gauge.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	fg := th.Color(theme.ColorNameForeground, v)

	r.face.FillColor = th.Color(theme.ColorNameInputBorder, v)
	r.center.FillColor = fg
	r.needle.StrokeColor = th.Color(theme.ColorNamePrimary, v)
	r.title.Color = fg
	r.readout.Color = fg
	for _, l := range r.labels {
		if l != nil {
			l.Color = fg
		}
	}
}

func (r *gaugeRenderer) formatValue() string {
	if f := r.gauge.TextFormatter; f != nil {
		return f()
	}
	return strconv.FormatFloat(r.gauge.Value, 'f', -1, 64)
}

// placeLine positions l along the direction given by sin/cos, starting offset
// away from the dial centre and extending length further out.
func (r *gaugeRenderer) placeLine(l *canvas.Line, sin, cos, offset, length float32) {
	x := r.middle.X + offset*sin
	y := r.middle.Y - offset*cos
	l.Position1 = fyne.NewPos(x, y)
	l.Position2 = fyne.NewPos(x+length*sin, y-length*cos)
}

func (r *gaugeRenderer) rotateNeedle() {
	g := r.gauge
	span := g.Max - g.Min
	if span <= 0 {
		span = 1
	}
	v := g.Value
	if v < g.Min {
		v = g.Min
	}
	if v > g.Max {
		v = g.Max
	}
	sin, cos := math.Sincos(gaugeStart + gaugeSweep*(v-g.Min)/span)
	r.placeLine(r.needle, float32(sin), float32(cos), r.needleOffset, r.needleLength)
}

func (r *gaugeRenderer) Layout(size fyne.Size) {
	diameter := fyne.Min(size.Width, size.Height)
	radius := diameter / 2
	r.middle = fyne.NewPos(size.Width/2, size.Height/2)
	r.needleOffset = -radius * 0.15
	r.needleLength = radius * 1.14

	r.face.Move(r.middle.SubtractXY(radius, radius))
	r.face.Resize(fyne.NewSquareSize(diameter))

	center := radius / 4
	r.center.Move(r.middle.SubtractXY(center/2, center/2))
	r.center.Resize(fyne.NewSquareSize(center))

	r.title.TextSize = radius / 4
	r.title.Move(r.middle.AddXY(0, diameter/3.5))

	r.readout.TextSize = radius / 3
	r.readout.Move(r.middle.SubtractXY(radius, radius).AddXY(0, diameter/5))
	r.readout.Resize(fyne.NewSquareSize(diameter))

	r.needle.StrokeWidth = diameter / 60
	r.rotateNeedle()

	majorStroke := fyne.Max(2, diameter/80)
	minorStroke := fyne.Max(2, diameter/200)
	labelPad := fyne.Max(6, radius*0.14)
	labelTextSize := radius * 0.1
	for i, graduation := range r.graduations {
		sin64, cos64 := math.Sincos(gaugeStart + gaugeSweep*float64(i)/float64(r.graduationCount))
		sin, cos := float32(sin64), float32(cos64)

		lbl := r.labels[i]
		if lbl == nil { // minor graduation
			graduation.StrokeWidth = minorStroke
			r.placeLine(graduation, sin, cos, radius*0.875, radius*0.125-2)
			continue
		}

		graduation.StrokeWidth = majorStroke
		r.placeLine(graduation, sin, cos, radius*0.75, radius*0.25-2)

		lbl.TextSize = labelTextSize
		lblSize := fyne.MeasureText(lbl.Text, labelTextSize, lbl.TextStyle)
		lblRadius := radius*0.75 - labelPad // labels sit on the inside of the major graduations
		lbl.Resize(lblSize)
		lbl.Move(fyne.NewPos(r.middle.X+sin*lblRadius-lblSize.Width/2, r.middle.Y-cos*lblRadius-lblSize.Height/2))
	}
}

func (r *gaugeRenderer) MinSize() fyne.Size {
	return r.gauge.minSize
}

func (r *gaugeRenderer) Refresh() {
	g := r.gauge
	full := g.Min != r.builtMin || g.Max != r.builtMax || g.Graduations != r.builtGraduations
	if full {
		r.rebuildGraduations()
		r.Layout(g.Size())
	}

	fg := g.Theme().Color(theme.ColorNameForeground, fyne.CurrentApp().Settings().ThemeVariant())
	if full || fg != r.builtFg {
		r.builtFg = fg
		r.applyTheme()
		full = true
	}

	if g.Title != r.title.Text {
		r.title.Text = g.Title
		if !full {
			r.title.Refresh()
		}
	}

	r.rotateNeedle()
	if s := r.formatValue(); s != r.readout.Text {
		r.readout.Text = s
		if !full {
			r.readout.Refresh()
		}
	}

	if full {
		for _, o := range r.objects {
			o.Refresh()
		}
		return
	}
	r.needle.Refresh()
}

func (r *gaugeRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *gaugeRenderer) Destroy() {
}
