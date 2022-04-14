package charts

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// BarChart or BarChart (alias).
type BarChart struct {
	*LineChart // BarChart overrides LineChart rasterization.
}

// BarChartOptions aliased to LineCharthOpts
type BarChartOptions = LineCharthOpts

// NewBarChart returns a new BarChart.
func NewBarChart(opts *BarChartOptions) *BarChart {
	chart := new(BarChart)
	chart.LineChart = NewLineChart(opts)
	return chart
}

// CreateRenderer creates a simple renderer
func (chart *BarChart) CreateRenderer() fyne.WidgetRenderer {
	w := chart.LineChart.CreateRenderer()
	chart.LineChart.image = canvas.NewRaster(chart.rasterize) // force the HistrogramChart rasterizer

	// recreate the container
	chart.LineChart.canvas = container.NewWithoutLayout(
		chart.LineChart.image,
		chart.LineChart.overlay,
	)

	// change the first object, leave the "overlay" in place
	w.Objects()[0] = chart.LineChart.canvas

	return w
}

// rasterize is called by the "image" Raster object.
func (chart *BarChart) rasterize(w, h int) image.Image {
	if chart.canvas == nil || len(chart.data) == 0 {
		return image.NewAlpha(image.Rect(0, 0, w, h))
	}

	// <!> Force the width and height to be the same as the image size
	// To not do this will cause the graph to be scaled down.
	// TODO: why is this needed?
	w = int(chart.image.Size().Width)
	h = int(chart.image.Size().Height)

	// Calculate the max and min values to scale the graph
	// and the step on X to move for each "point"
	width := float64(w)
	height := float64(h)
	stepX := width / float64(len(chart.data))
	maxY := float64(0)
	minY := float64(0)
	if chart.Options().GraphRange == nil {
		for _, v := range chart.data {
			if v > maxY {
				maxY = v
			}
			if v < minY {
				minY = v
			}
		}
	} else {
		maxY = chart.Options().GraphRange.YMax
		minY = chart.Options().GraphRange.YMin
	}

	// reduction factor
	reduce := height / (maxY - minY)

	// keep the Y fix value - used by GetDataPosAt()
	chart.yFix = [2]float64{minY, reduce}

	// Draw...
	currentX := float64(0)

	// each "value" has 4 points (bottom left, top left, top right, bottom right)
	// each point is defined by 2 coordinates (x, y)
	points := make([][2]float64, len(chart.data)*4+1)

	sw := float64(chart.opts.StrokeWidth)

	for i, v := range chart.data {
		// Calculate the points
		// bottom left
		points[i*4+0][0] = currentX
		points[i*4+0][1] = height + sw
		// top left
		points[i*4+1][0] = currentX
		points[i*4+1][1] = height - (v-minY)*reduce + sw
		// top right
		points[i*4+2][0] = currentX + stepX
		points[i*4+2][1] = height - (v-minY)*reduce + sw
		// bottom right
		points[i*4+3][0] = currentX + stepX
		points[i*4+3][1] = height + sw

		currentX += stepX
	}

	points[len(points)-1][0] = currentX
	points[len(points)-1][1] = height

	// colors
	fgR, fgG, fgB, _ := chart.opts.StrokeColor.RGBA()
	bgR, bgG, bgB, _ := chart.opts.FillColor.RGBA()
	// convert the svg to an image.Image
	buff := new(bytes.Buffer)
	err := getLineSVGTemplate().Execute(buff, svgTplLineStruct{
		Data:        points,
		Width:       w,
		Height:      h,
		StrokeWidth: chart.opts.StrokeWidth,
		StrokeColor: fmt.Sprintf("#%02x%02x%02x", uint8(fgR/0x101), uint8(fgG/0x101), uint8(fgB/0x101)),
		Fill:        fmt.Sprintf("#%02x%02x%02x", uint8(bgR/0x101), uint8(bgG/0x101), uint8(bgB/0x101)),
	})

	if err != nil {
		log.Println(err)
		return image.NewAlpha(image.Rect(0, 0, w, h))
	}

	graph, err := oksvg.ReadIconStream(buff)
	if err != nil {
		log.Println(err)
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	graph.SetTarget(0, 0, float64(w), float64(h))
	scanner := rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())
	graph.Draw(rasterx.NewDasher(w, h, scanner), 1)

	return rgba
}
