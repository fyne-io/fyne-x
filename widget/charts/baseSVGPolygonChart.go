package charts

import (
	"image/color"
	"sync"
	"text/template"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const svgLineTplString = `<svg xmlns="http://www.w3.org/2000/svg" width="{{ .Width }}" height="{{ .Height }}" viewBox="0 0 {{.Width}} {{ .Height }}">
    <polygon 
        points="{{ range .Data }}{{ index . 0 }},{{ index . 1 }} {{ end }}"
        style="{{ if ne .FillColor "" }}fill:{{.FillColor}};{{end}}{{ if ne .StrokeColor "" }}stroke:{{ .StrokeColor }};{{ end }}stroke-width:{{ .StrokeWidth }}"
    />
</svg>`

// svgPolygonTpl is the template.Templatethat can be used by the SVG renderers.
var svgPolygonTpl *template.Template

// GetPolygonSVGTemplate return initialized and executed template.Template for polygon.
func GetPolygonSVGTemplate() *template.Template {
	if svgPolygonTpl == nil {
		svgPolygonTpl = template.Must(template.New("polygon").Parse(svgLineTplString))
	}
	return svgPolygonTpl
}

// SVGTplPolygonStruct  handles the graph data, colors... for Line SVG
type SVGTplPolygonStruct struct {
	// Width of the chart
	Width int

	// Height of the chart
	Height int

	// Data points (X,Y) to draw the SVG polygon
	Data [][2]float64

	// Fill color in SVG/HTML format (e.g. #ff0000, red, none, ...)
	FillColor string

	// Stroke color in SVG/HTML format (e.g. #ff0000, red, none, ...)
	StrokeColor string

	// Stroke width of the line
	StrokeWidth float32
}

// PolygonCharthOpts provides options for the polygon chart.
type PolygonCharthOpts struct {
	*sizeOpts

	// FillColor is the color of the fill. Alpha is ignored.
	FillColor color.Color

	// StrokeWidth is the width of the stroke.
	StrokeWidth float32

	// StrokeColor is the color of the stroke. Alpha is ignored.
	StrokeColor color.Color
}

// BasePolygonSVGChart is the base widget to implement new chart widget with SVG using a polygon element. This should be not use directly but used to create a new chart.
type BasePolygonSVGChart struct {
	widget.BaseWidget
	*BaseSVGChart
	opts    *PolygonCharthOpts
	canvas  *fyne.Container
	overlay *fyne.Container
	locker  *sync.Mutex

	yFix [2]float64
	data []float64
}

// NewBasePolygonSVGChart creates a new BasePolygonSVGChart.
func NewBasePolygonSVGChart(options *PolygonCharthOpts) *BasePolygonSVGChart {
	g := &BasePolygonSVGChart{
		locker: &sync.Mutex{},
	}
	g.BaseSVGChart = NewSVGGraph()
	if options != nil {
		g.opts = options
	} else {
		g.opts = &PolygonCharthOpts{
			sizeOpts:    newSizeOpts(),
			StrokeWidth: 1,
			StrokeColor: theme.ForegroundColor(),
			FillColor:   theme.DisabledButtonColor(),
		}
	}

	if g.opts.StrokeColor == nil {
		g.opts.StrokeColor = theme.ForegroundColor()
	}

	if g.opts.StrokeWidth == 0 {
		g.opts.StrokeWidth = .7
	}

	if g.opts.FillColor == nil {
		g.opts.FillColor = theme.DisabledButtonColor()
	}
	return g
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (g *BasePolygonSVGChart) CreateRenderer() fyne.WidgetRenderer {
	g.overlay = container.NewWithoutLayout()
	g.canvas = container.NewWithoutLayout(g.Rasterizer(), g.overlay)
	return widget.NewSimpleRenderer(g.canvas)
}

// GetDrawable returns the graph's overlay drawable container.
func (g *BasePolygonSVGChart) GetDrawable() *fyne.Container {
	return g.overlay
}

// CalculateGraphScale calculates the scale of the graph. It returns Ymin, YMax, stepX and reduce factor.
func (g *BasePolygonSVGChart) CalculateGraphScale(w, h int) (float64, float64, float64, float64) {

	var (
		maxY float64
		minY float64
	)
	width := float64(w)
	height := float64(h)
	stepX := width / float64(len(g.data))
	if g.Options().GraphRange == nil {
		for _, v := range g.data {
			if v > maxY {
				maxY = v
			}
			if v < minY {
				minY = v
			}
		}
	} else {
		maxY = g.Options().GraphRange.YMax
		minY = g.Options().GraphRange.YMin
	}

	// reduction factor
	reduce := height / (maxY - minY)

	// keep the Y fix value - used by GetDataPosAt()
	g.yFix = [2]float64{minY, reduce}
	return minY, maxY, stepX, reduce
}

// MinSize returns the smallest size this widget can shrink to.
func (g *BasePolygonSVGChart) MinSize() fyne.Size {
	return g.BaseWidget.MinSize()
}

// GetDataPosAt returns the data value and and the exact position on the curve for a given position. This is
// useful to draw something on the graph at mouse position for example.
func (g *BasePolygonSVGChart) GetDataPosAt(pos fyne.Position) (float64, fyne.Position) {

	if len(g.data) == 0 {
		return 0, fyne.NewPos(0, 0)
	}

	stepX := g.Rasterizer().Size().Width / float32(len(g.data))
	// get the X value corresponding to the data index
	x := int(pos.X / g.Rasterizer().Size().Width * float32(len(g.data)))
	if x < 0 || x >= len(g.data) {
		return 0, fyne.NewPos(0, 0)
	}
	value := g.data[int(x)]

	// now, get the Y value corresponding to the data value
	y := float64(g.Rasterizer().Size().Height) - value*g.yFix[1] + g.yFix[0]*g.yFix[1]

	// calculate the X value on the graph
	xp := float32(x) * stepX

	return value, fyne.NewPos(xp, float32(y))
}

// Options returns the options of the graph. You can change the options after the graph is created.
func (g *BasePolygonSVGChart) Options() *PolygonCharthOpts {
	if g.opts == nil {
		g.opts = &PolygonCharthOpts{}
	}
	return g.opts
}

// SetData sets the data to be displayed in the graph
func (g *BasePolygonSVGChart) SetData(data []float64) {
	g.locker.Lock()
	defer g.locker.Unlock()
	g.data = data
}

// Size returns the size of the graph widget.
func (g *BasePolygonSVGChart) Size() fyne.Size {
	if g.canvas == nil {
		return fyne.NewSize(0, 0)
	}
	return g.canvas.Size()
}

// Resize sets a new size for the graph.
func (g *BasePolygonSVGChart) Resize(size fyne.Size) {
	g.BaseWidget.Resize(size)
	if g.canvas != nil {
		g.canvas.Resize(size)
		g.Rasterizer().Resize(size)
		g.overlay.Resize(size)
	}
	g.Refresh()
}

// Refresh refreshes the graph.
func (g *BasePolygonSVGChart) Refresh() {

	if g.canvas == nil {
		return
	}
	g.BaseWidget.Refresh()
	for _, child := range g.overlay.Objects {
		child.Refresh()
	}
	g.Rasterizer().Refresh()
}
