package charts

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// LineCharthOpts provides options for the graph.
type LineCharthOpts struct {
	*globalOpts

	// FillColor is the color of the fill. Alpha is ignored.
	FillColor color.Color

	// StrokeWidth is the width of the stroke.
	StrokeWidth float32

	// StrokeColor is the color of the stroke. Alpha is ignored.
	StrokeColor color.Color

	// Title is the title of the graph.
}

// LineChart widget provides a plotting widget for data.
type LineChart struct {
	graph
	canvas  *fyne.Container
	overlay *fyne.Container
	data    []float64
	image   *canvas.Raster
	locker  sync.Mutex
	opts    *LineCharthOpts
	yFix    [2]float64
}

// NewLineChart creates a new graph widget. The "options" parameter is optional. IF you provide several options, only the first will be used.
func NewLineChart(options *LineCharthOpts) *LineChart {
	g := &LineChart{
		data:   []float64{},
		locker: sync.Mutex{},
		yFix:   [2]float64{},
	}

	if options != nil {
		g.opts = options
	} else {
		g.opts = &LineCharthOpts{
			globalOpts:  newGlobalOpts(),
			StrokeWidth: 1,
			StrokeColor: theme.ForegroundColor(),
			FillColor:   theme.DisabledButtonColor(),
		}
	}

	if g.opts.StrokeColor == nil {
		g.opts.StrokeColor = theme.ForegroundColor()
	}

	if g.opts.StrokeWidth == 0 {
		g.opts.StrokeWidth = .7
	}

	if g.opts.FillColor == nil {
		g.opts.FillColor = theme.DisabledButtonColor()
	}

	g.ExtendBaseWidget(g)

	return g
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (g *LineChart) CreateRenderer() fyne.WidgetRenderer {
	g.image = canvas.NewRaster(g.rasterize)
	g.overlay = container.NewWithoutLayout()
	g.canvas = container.NewWithoutLayout(g.image, g.overlay)
	return widget.NewSimpleRenderer(g.canvas)
}

// GetDrawable returns the graph's overlay drawable container.
func (g *LineChart) GetDrawable() *fyne.Container {
	return g.overlay
}

// GetDataPosAt returns the data value and and the exact position on the curve for a given position. This is
// useful to draw something on the graph at mouse position for example.
func (g *LineChart) GetDataPosAt(pos fyne.Position) (float64, fyne.Position) {

	if len(g.data) == 0 {
		return 0, fyne.NewPos(0, 0)
	}

	if g.image == nil {
		return 0, fyne.NewPos(0, 0)
	}

	stepX := g.image.Size().Width / float32(len(g.data))
	// get the X value corresponding to the data index
	x := int(pos.X / g.image.Size().Width * float32(len(g.data)))
	if x < 0 || x >= len(g.data) {
		return 0, fyne.NewPos(0, 0)
	}
	value := g.data[int(x)]

	// now, get the Y value corresponding to the data value
	y := float64(g.image.Size().Height) - value*g.yFix[1] + g.yFix[0]*g.yFix[1]

	// calculate the X value on the graph
	xp := float32(x) * stepX

	return value, fyne.NewPos(xp, float32(y))
}

// MinSize returns the smallest size this widget can shrink to.
func (g *LineChart) MinSize() fyne.Size {
	if g.image == nil {
		return fyne.NewSize(0, 0)
	}
	return g.BaseWidget.MinSize()
}

// Options returns the options of the graph. You can change the options after the graph is created.
func (g *LineChart) Options() *LineCharthOpts {
	if g.opts == nil {
		g.opts = &LineCharthOpts{}
	}
	return g.opts
}

// Refresh refreshes the graph.
func (g *LineChart) Refresh() {

	if g.image == nil {
		return
	}
	g.image.Refresh()
	g.canvas.Refresh()
}

// Resize sets a new size for the graph.
func (g *LineChart) Resize(size fyne.Size) {
	g.BaseWidget.Resize(size)
	if g.canvas != nil {
		g.canvas.Resize(size)
		g.image.Resize(size)
		g.overlay.Resize(size)
	}
	g.Refresh()
}

// SetData sets the data for the graph - each call to this method will redraw the graph.
func (g *LineChart) SetData(data []float64) {
	g.locker.Lock()
	g.data = data
	g.locker.Unlock()
	g.Refresh()
}

// Size returns the size of the graph widget.
func (g *LineChart) Size() fyne.Size {
	if g.canvas == nil {
		return fyne.NewSize(0, 0)
	}
	return g.canvas.Size()
}

// This private method is linjed to g.image canvas.Raster property. It uses oksvg and rasterx to render the graph from a SVG template.
func (g *LineChart) rasterize(w, h int) image.Image {

	g.locker.Lock()
	defer g.locker.Unlock()

	if g.image == nil || len(g.data) == 0 {
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}

	// <!> Force the width and height to be the same as the image size
	// To not do this will cause the graph to be scaled down.
	// TODO: why is this needed?
	w = int(g.image.Size().Width)
	h = int(g.image.Size().Height)

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
	err := getLineSVGTemplate().Execute(buff, svgTplLineStruct{
		Data:        points,
		Width:       w,
		Height:      h,
		StrokeWidth: g.opts.StrokeWidth,
		StrokeColor: fmt.Sprintf("#%02x%02x%02x", uint8(fgR/0x101), uint8(fgG/0x101), uint8(fgB/0x101)),
		Fill:        fmt.Sprintf("#%02x%02x%02x", uint8(bgR/0x101), uint8(bgG/0x101), uint8(bgB/0x101)),
	})
	if err != nil {
		log.Println(err)
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}

	// convert the svg to an image.Image
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
