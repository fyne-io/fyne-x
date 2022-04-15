package charts

import (
	"image"
	"io"
	"log"

	"fyne.io/fyne/v2/canvas"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
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

// base struct for any Graph object.
type BaseChart struct{}

// BaseSVGChart uses SVG to render the graph.
type BaseSVGChart struct {
	*BaseChart

	rasterizer *canvas.Raster
}

// NewSVGGraph returns a new SVG graph.
func NewSVGGraph() *BaseSVGChart {
	g := &BaseSVGChart{}
	return g
}

// SetRasterizerFunc sets the rasterizer function for the graph rasterizer.
func (s *BaseSVGChart) SetRasterizerFunc(r func(w, h int) image.Image) {
	s.Rasterizer().Generator = r
}

// Rasterizer returns the rasterizer for the graph.
func (s *BaseSVGChart) Rasterizer() *canvas.Raster {
	if s.rasterizer == nil {
		s.rasterizer = canvas.NewRaster(func(w, h int) image.Image {
			return image.NewRGBA(image.Rect(0, 0, w, h))
		})
	}
	return s.rasterizer
}

// Render the SVG string from buff to an image.Imaage.
func (g *BaseSVGChart) Render(buff io.Reader, width, height int) image.Image {

	graph, err := oksvg.ReadIconStream(buff)
	if err != nil {
		log.Println(err)
		return image.NewRGBA(image.Rect(0, 0, width, height))
	}
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	graph.SetTarget(0, 0, float64(width), float64(height))
	scanner := rasterx.NewScannerGV(width, height, rgba, rgba.Bounds())
	graph.Draw(rasterx.NewDasher(width, height, scanner), 1)
	return rgba
}
