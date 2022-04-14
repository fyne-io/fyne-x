package charts

import (
	"fmt"
	"text/template"

	"fyne.io/fyne/v2/widget"
)

const svgLineTplString = `<svg xmlns="http://www.w3.org/2000/svg" width="{{.Width}}" height="{{.Height}}" viewBox="0 0 {{.Width}} {{.Height}}">
    <polygon 
        points="{{range .Data}}{{index . 0}},{{ index . 1}} {{end}}"
        style="fill:{{.Fill}};stroke:{{.StrokeColor}};stroke-width:{{.StrokeWidth}}"
    />
</svg>`

var svgLineTpl *template.Template

// GraphRange set the range of the graph.
type GraphRange struct {
	// Set the Y range of the graph.
	YMin, YMax float64

	// Set the X range of the graph.
	// XMin, XMax float64
}

// base struct for any Graph object.
type graph struct {
	widget.BaseWidget
	graphRange *GraphRange
}

// SetGraphRange set the GraphRange of the chart. If the range is nil, so the
// chart will generallly auto scale the range.
func (g *graph) SetGraphRange(gr *GraphRange) error {

	if gr == nil {
		return nil
	}

	// check that min < max
	if gr.YMin >= gr.YMax {
		return fmt.Errorf("YMin must be less than YMax")
	}

	g.graphRange = gr
	return nil
}

// return initialized svgLineTpl
func getLineSVGTemplate() *template.Template {
	if svgLineTpl == nil {
		svgLineTpl = template.Must(template.New("svg").Parse(svgLineTplString))
	}
	return svgLineTpl
}

// structure to handle the graph data, colors... for Line SVG
type svgTplLineStruct struct {
	Width       int
	Height      int
	Data        [][2]float64
	Fill        string
	StrokeColor string
	StrokeWidth float32
}
