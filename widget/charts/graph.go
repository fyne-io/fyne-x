package charts

import (
	"fyne.io/fyne/v2/canvas"
)

// GraphRange set the range of the graph.
type GraphRange struct {
	// Set the Y range of the graph.
	YMin, YMax float64
}

// global options for the chart
type sizeOpts struct {
	GraphRange *GraphRange
}

func newSizeOpts() *sizeOpts {
	return &sizeOpts{}
}

// SetGraphRange set the global range of the graph.
func (g *sizeOpts) SetGraphRange(r *GraphRange) {
	g.GraphRange = r
}

// BaseChart struct for any Graph object.
type BaseChart struct{}

// BaseSVGChart uses SVG to render the graph.
type BaseSVGChart struct {
	*BaseChart

	rasterizer *canvas.Raster
}
