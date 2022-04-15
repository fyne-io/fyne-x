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

func (chart *LineChart) rasterize(w, h int) image.Image {
	chart.Lock()
	defer chart.Unlock()

	if len(chart.data) == 0 {
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}

	// <!> Force the width and height to be the same as the image size
	// To not do this will cause the graph to be scaled down.
	// TODO: why is this needed?
	w = int(chart.Rasterizer().Size().Width)
	h = int(chart.Rasterizer().Size().Height)

	// prepare points
	points := make([][2]float64, len(chart.data)+3)

	// colors
	strokeColor := "none"
	fillColor := "none"
	transp := color.RGBA{0, 0, 0, 0}
	if chart.opts.StrokeColor != transp {
		fgR, fgG, fgB, _ := chart.opts.StrokeColor.RGBA()
		strokeColor = fmt.Sprintf("#%02x%02x%02x", uint8(fgR/0x101), uint8(fgG/0x101), uint8(fgB/0x101))
	}
	if chart.opts.FillColor != transp {
		bgR, bgG, bgB, _ := chart.opts.FillColor.RGBA()
		fillColor = fmt.Sprintf("#%02x%02x%02x", uint8(bgR/0x101), uint8(bgG/0x101), uint8(bgB/0x101))
	}

	// Calculate the max and min values to scale the graph
	// and the step on X to move for each "point"
	width := float64(w)
	height := float64(h)

	// build point positions
	var currentX float64
	minY, _, stepX, reduce := chart.GraphScale(w, h)
	strokeWidth := float64(chart.opts.StrokeWidth)
	// place some start and end points to make the graph look nicer
	points[0] = [2]float64{strokeWidth, height + minY*reduce + strokeWidth}
	points[len(points)-2] = [2]float64{
		width + strokeWidth,
		height - chart.data[len(chart.data)-1]*reduce + minY*reduce,
	}
	points[len(points)-1] = [2]float64{
		width + strokeWidth*2,
		height + minY*reduce + strokeWidth,
	}

	// add points from data
	for i, d := range chart.data {
		y := float64(0)
		if d < 0 {
			y = height - strokeWidth - d*reduce + minY*reduce
		} else if d > 0 {
			y = height + strokeWidth - d*reduce + minY*reduce
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
		StrokeWidth: chart.opts.StrokeWidth,
		StrokeColor: strokeColor,
		FillColor:   fillColor,
	})
	if err != nil {
		log.Println(err)
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}

	return chart.Render(buff, w, h)
}
