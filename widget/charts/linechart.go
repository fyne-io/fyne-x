package charts

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
)

// LineChart widget provides a plotting widget for data.
type LineChart struct {
	*BasePolygonSVGChart
}

// NewLineChart creates a new graph widget. The "options" parameter is optional. IF you provide several options, only the first will be used.
func NewLineChart(options *PolygonCharthOpts) *LineChart {
	g := &LineChart{
		BasePolygonSVGChart: NewBasePolygonSVGChart(options),
	}

	g.ExtendBaseWidget(g)
	g.SetRasterizerFunc(g.rasterize)

	return g
}

func (g *LineChart) rasterize(w, h int) image.Image {

	g.locker.Lock()
	defer g.locker.Unlock()

	if len(g.data) == 0 {
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}

	// <!> Force the width and height to be the same as the image size
	// To not do this will cause the graph to be scaled down.
	// TODO: why is this needed?
	w = int(g.Rasterizer().Size().Width)
	h = int(g.Rasterizer().Size().Height)

	// prepare points
	points := make([][2]float64, len(g.data)+3)

	// colors
	strokeColor := "none"
	fillColor := "none"
	transp := color.RGBA{0, 0, 0, 0}
	if g.opts.StrokeColor != transp {
		fgR, fgG, fgB, _ := g.opts.StrokeColor.RGBA()
		strokeColor = fmt.Sprintf("#%02x%02x%02x", uint8(fgR/0x101), uint8(fgG/0x101), uint8(fgB/0x101))
	}
	if g.opts.FillColor != transp {
		bgR, bgG, bgB, _ := g.opts.FillColor.RGBA()
		fillColor = fmt.Sprintf("#%02x%02x%02x", uint8(bgR/0x101), uint8(bgG/0x101), uint8(bgB/0x101))
	}

	// Calculate the max and min values to scale the graph
	// and the step on X to move for each "point"
	width := float64(w)
	height := float64(h)

	// build point positions
	var currentX float64
	minY, _, stepX, reduce := g.CalculateGraphScale(w, h)
	sw := float64(g.opts.StrokeWidth)
	points[0] = [2]float64{sw, height + minY*reduce + sw}
	points[len(points)-2] = [2]float64{width + sw, height - g.data[len(g.data)-1]*reduce + minY*reduce}
	points[len(points)-1] = [2]float64{width + sw*2, height + minY*reduce + sw}

	for i, d := range g.data {
		y := float64(0)
		if d < 0 {
			y = height - sw - d*reduce + minY*reduce
		} else if d > 0 {
			y = height + sw - d*reduce + minY*reduce
		} else {
			y = height + minY*reduce
		}
		points[i+1] = [2]float64{currentX, y}
		currentX += stepX
	}

	// render SVG template
	buff := new(bytes.Buffer)
	err := GetPolygonSVGTemplate().Execute(buff, SVGTplPolygonStruct{
		Data:        points,
		Width:       w,
		Height:      h,
		StrokeWidth: g.opts.StrokeWidth,
		StrokeColor: strokeColor,
		FillColor:   fillColor,
	})
	if err != nil {
		log.Println(err)
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}

	return g.Render(buff, w, h)
}
