package charts

import (
	"image"
	"io"

	"fyne.io/fyne/v2/canvas"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// BaseSVGChart uses SVG to render the graph.
type BaseSVGChart struct {
	*BaseChart

	rasterizer *canvas.Raster
}

// NewSVGGraph returns a new SVG graph.
func NewSVGGraph() *BaseSVGChart {
	g := &BaseSVGChart{
		BaseChart: &BaseChart{},
	}
	return g
}

// Rasterizer returns the rasterizer for the graph.
func (chart *BaseSVGChart) Rasterizer() *canvas.Raster {
	if chart.rasterizer == nil {
		chart.rasterizer = canvas.NewRaster(func(w, h int) image.Image {
			return image.NewRGBA(image.Rect(0, 0, w, h))
		})
	}
	return chart.rasterizer
}

// Render the SVG string from buff to an image.Imaage.
func (chart *BaseSVGChart) Render(buff io.Reader, width, height int) image.Image {
	graph, _ := oksvg.ReadIconStream(buff)
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	graph.SetTarget(0, 0, float64(width), float64(height))
	scanner := rasterx.NewScannerGV(width, height, rgba, rgba.Bounds())
	graph.Draw(rasterx.NewDasher(width, height, scanner), 1)
	return rgba
}

// SetRasterizerFunc sets the rasterizer function for the graph rasterizer.
func (chart *BaseSVGChart) SetRasterizerFunc(r func(w, h int) image.Image) {
	chart.Rasterizer().Generator = r
}
