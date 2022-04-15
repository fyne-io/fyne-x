package charts

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
)

// BarChart or BarChart (alias).
type BarChart struct {
	*BasePolygonSVGChart
}

// BarChartOptions aliased to LineCharthOpts
type BarChartOptions = PolygonCharthOpts

// NewBarChart returns a new BarChart.
func NewBarChart(opts *BarChartOptions) *BarChart {
	chart := &BarChart{
		BasePolygonSVGChart: NewBasePolygonSVGChart(opts),
	}
	chart.SetRasterizerFunc(chart.rasterize)

	chart.ExtendBaseWidget(chart)
	return chart
}

// rasterize is called by the "image" Raster object.
func (chart *BarChart) rasterize(w, h int) image.Image {
	chart.Lock()
	defer chart.Unlock()
	if chart.canvas == nil || len(chart.data) == 0 {
		return image.NewAlpha(image.Rect(0, 0, w, h))
	}

	// <!> Force the width and height to be the same as the image size
	// To not do this will cause the graph to be scaled down.
	// TODO: why is this needed?
	w = int(chart.Rasterizer().Size().Width)
	h = int(chart.Rasterizer().Size().Height)

	// Calculate the max and min values to scale the graph
	// and the step on X to move for each "point"
	height := float64(h)

	var currentX float64
	minY, _, stepX, reduce := chart.GraphScale(w, h)

	// each "value" has 4 points (bottom left, top left, top right, bottom right)
	// each point is defined by 2 coordinates (x, y)
	points := make([][2]float64, len(chart.data)*4+1)

	strokeWidth := float64(chart.opts.StrokeWidth)

	for i, v := range chart.data {
		// Calculate the points
		// bottom left
		points[i*4+0][0] = currentX
		points[i*4+0][1] = height + strokeWidth
		// top left
		points[i*4+1][0] = currentX
		points[i*4+1][1] = height - (v-minY)*reduce + strokeWidth
		// top right
		points[i*4+2][0] = currentX + stepX
		points[i*4+2][1] = height - (v-minY)*reduce + strokeWidth
		// bottom right
		points[i*4+3][0] = currentX + stepX
		points[i*4+3][1] = height + strokeWidth

		currentX += stepX
	}

	points[len(points)-1][0] = currentX
	points[len(points)-1][1] = height

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
	// convert the svg to an image.Image
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
		return image.NewAlpha(image.Rect(0, 0, w, h))
	}

	return chart.Render(buff, w, h)
}
