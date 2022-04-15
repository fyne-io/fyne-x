package charts

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"sync"

	"fyne.io/fyne/v2"
)

// LineChart widget provides a plotting widget for data.
type LineChart struct {
	*BasePolygonSVGChart
	data   []float64
	locker sync.Mutex
	yFix   [2]float64
}

// NewLineChart creates a new graph widget. The "options" parameter is optional. IF you provide several options, only the first will be used.
func NewLineChart(options *PolygonCharthOpts) *LineChart {
	g := &LineChart{
		BasePolygonSVGChart: newPolygonChart(options),
		data:                []float64{},
		locker:              sync.Mutex{},
		yFix:                [2]float64{},
	}

	g.ExtendBaseWidget(g)
	g.SetRasterizerFunc(g.rasterize)

	return g
}

// GetDataPosAt returns the data value and and the exact position on the curve for a given position. This is
// useful to draw something on the graph at mouse position for example.
func (g *LineChart) GetDataPosAt(pos fyne.Position) (float64, fyne.Position) {

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

// SetData sets the data for the graph - each call to this method will redraw the graph.
func (g *LineChart) SetData(data []float64) {
	g.locker.Lock()
	g.data = data
	g.locker.Unlock()
	g.Refresh()
}

// This private method is linjed to g.image canvas.Raster property. It uses oksvg and rasterx to render the graph from a SVG template.
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
	fgR, fgG, fgB, _ := g.opts.StrokeColor.RGBA()
	bgR, bgG, bgB, _ := g.opts.FillColor.RGBA()

	// Calculate the max and min values to scale the graph
	// and the step on X to move for each "point"
	width := float64(w)
	height := float64(h)
	stepX := width / float64(len(g.data))
	maxY := float64(0)
	minY := float64(0)
	if g.Options().GraphRange == nil {
		for _, v := range g.data {
			if v > maxY {
				maxY = v
			}
			if v < minY {
				minY = v
			}
		}

		// Move the graph to avoid the "zero" line
		if minY > 0 {
			minY = 0
		}
	} else {
		maxY = g.Options().GraphRange.YMax
		minY = g.Options().GraphRange.YMin
	}

	// reduction factor
	reduce := height / (maxY - minY)

	// keep the Y fix value - used by GetDataPosAt()
	g.yFix = [2]float64{minY, reduce}

	// build point positions
	currentX := float64(0)
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
	err := getPolygonSVGTemplate().Execute(buff, svgTplLineStruct{
		Data:        points,
		Width:       w,
		Height:      h,
		StrokeWidth: g.opts.StrokeWidth,
		StrokeColor: fmt.Sprintf("#%02x%02x%02x", uint8(fgR/0x101), uint8(fgG/0x101), uint8(fgB/0x101)),
		FillColor:   fmt.Sprintf("#%02x%02x%02x", uint8(bgR/0x101), uint8(bgG/0x101), uint8(bgB/0x101)),
	})
	if err != nil {
		log.Println(err)
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}

	return g.Render(buff, w, h)
}
