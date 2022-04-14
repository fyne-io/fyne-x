package charts

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func createTestBarChart() *BarChart {
	return NewBarChart(nil)
}

func createTestBarChartWithOptions() *BarChart {
	return NewBarChart(&LineCharthOpts{
		StrokeColor: color.RGBA{0x11, 0x22, 0x33, 255},
		FillColor:   color.RGBA{0x44, 0x55, 0x66, 255},
		StrokeWidth: 5,
	})
}

func TestBarChart_Create(t *testing.T) {
	histo := createTestBarChart()
	assert.Equal(t, histo.opts.StrokeWidth, float32(1))
	assert.Equal(t, histo.opts.FillColor, theme.DisabledButtonColor())
	assert.Equal(t, histo.opts.StrokeColor, theme.ForegroundColor())

	assert.Len(t, histo.data, 0)
}

func TestBarChart_CreateWithOptions(t *testing.T) {
	graph := createTestBarChartWithOptions()
	assert.Equal(t, graph.opts.StrokeWidth, float32(5))
	assert.Equal(t, graph.opts.StrokeColor, color.RGBA{0x11, 0x22, 0x33, 255})
	assert.Equal(t, graph.opts.FillColor, color.RGBA{0x44, 0x55, 0x66, 255})
}

func TestBarChart_Rasterize(t *testing.T) {
	graph := createTestBarChart()

	win := test.NewWindow(graph)
	defer win.Close()
	win.Resize(fyne.NewSize(400, 200))

	data := []float64{1, 2, 3, 4, 5}
	graph.SetData(data)

	img := makeRasterize(win, graph)

	assertSize(t, img, graph)

}

func TestBarChart_GetOpts(t *testing.T) {
	opts := &BarChartOptions{
		StrokeWidth: 5,
		StrokeColor: color.RGBA{0x11, 0x22, 0x33, 255},
		FillColor:   color.RGBA{0x44, 0x55, 0x66, 255},
	}

	graph := NewBarChart(opts)
	assert.Equal(t, graph.opts.StrokeWidth, opts.StrokeWidth)
	assert.Equal(t, graph.opts.StrokeColor, opts.StrokeColor)
	assert.Equal(t, graph.opts.FillColor, opts.FillColor)

}

func TestBarChart_CreateRender(t *testing.T) {
	graph := createTestBarChart()

	win := test.NewWindow(graph)
	defer win.Close()
	win.Resize(fyne.NewSize(400, 200))

	r := graph.CreateRenderer()
	assert.NotNil(t, r)

	assert.NotNil(t, graph.canvas)
	assert.NotNil(t, graph.image)
	assert.NotNil(t, graph.overlay)

	assert.Len(t, graph.canvas.Objects, 2)

}
