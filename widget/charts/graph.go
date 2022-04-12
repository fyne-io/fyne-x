package charts

import (
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

// base struct for any Graph object.
type graph struct {
	widget.BaseWidget
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
