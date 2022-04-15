package charts

import (
	"image/color"
	"text/template"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const svgLineTplString = `<svg xmlns="http://www.w3.org/2000/svg" width="{{.Width}}" height="{{.Height}}" viewBox="0 0 {{.Width}} {{.Height}}">
    <polygon 
        points="{{range .Data}}{{index . 0}},{{ index . 1}} {{end}}"
        style="fill:{{.FillColor}};stroke:{{.StrokeColor}};stroke-width:{{.StrokeWidth}}"
    />
</svg>`

// svgPolygonTpl is the template.Templatethat can be used by the SVG renderers.
var svgPolygonTpl *template.Template

// return initialized svgLineTpl
func getPolygonSVGTemplate() *template.Template {
	if svgPolygonTpl == nil {
		svgPolygonTpl = template.Must(template.New("svg").Parse(svgLineTplString))
	}
	return svgPolygonTpl
}

// structure to handle the graph data, colors... for Line SVG
type svgTplLineStruct struct {
	Width       int
	Height      int
	Data        [][2]float64
	FillColor   string
	StrokeColor string
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

	// Title is the title of the graph.
}

// BasePolygonSVGChart is the base widget to implement new chart widget with SVG using a polygon element. This should be not use directly but used to create a new chart.
type BasePolygonSVGChart struct {
	widget.BaseWidget
	*BaseSVGChart
	opts    *PolygonCharthOpts
	canvas  *fyne.Container
	overlay *fyne.Container
}

func newPolygonChart(options *PolygonCharthOpts) *BasePolygonSVGChart {
	g := &BasePolygonSVGChart{}
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

func (g *BasePolygonSVGChart) CreateRenderer() fyne.WidgetRenderer {
	g.overlay = container.NewWithoutLayout()
	g.canvas = container.NewWithoutLayout(g.Rasterizer(), g.overlay)
	return widget.NewSimpleRenderer(g.canvas)
}

// GetDrawable returns the graph's overlay drawable container.
func (g *BasePolygonSVGChart) GetDrawable() *fyne.Container {
	return g.overlay
}

// MinSize returns the smallest size this widget can shrink to.
func (g *BasePolygonSVGChart) MinSize() fyne.Size {
	return g.BaseWidget.MinSize()
}

// Options returns the options of the graph. You can change the options after the graph is created.
func (g *BasePolygonSVGChart) Options() *PolygonCharthOpts {
	if g.opts == nil {
		g.opts = &PolygonCharthOpts{}
	}
	return g.opts
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
	g.Rasterizer().Refresh()
	g.canvas.Refresh()
}
