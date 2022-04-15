package charts

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func createTestLineChart() *LineChart {
	return NewLineChart(nil)
}

func createTestLineChartWithOptions() *LineChart {
	return NewLineChart(&PolygonCharthOpts{
		StrokeColor: color.RGBA{0x11, 0x22, 0x33, 255},
		FillColor:   color.RGBA{0x44, 0x55, 0x66, 255},
		StrokeWidth: 5,
	})
}

func TestLineChart_Creation(t *testing.T) {
	graph := createTestLineChart()

	assert.Len(t, graph.data, 0)
	assert.Equal(t, graph.opts.StrokeColor, theme.ForegroundColor())
	assert.Equal(t, graph.opts.FillColor, theme.DisabledButtonColor())
	assert.Equal(t, graph.opts.StrokeWidth, float32(1))
}

func TestLineChart_CreationWithOptions(t *testing.T) {
	graph := createTestLineChartWithOptions()

	assert.Equal(t, graph.opts.StrokeColor, color.RGBA{0x11, 0x22, 0x33, 255})
	assert.Equal(t, graph.opts.FillColor, color.RGBA{0x44, 0x55, 0x66, 255})
	assert.Equal(t, graph.opts.StrokeWidth, float32(5))
}

func TestLineChart_Rasterizer(t *testing.T) {
	graph := createTestLineChart()
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	graph.SetData(data)
	win := test.NewWindow(graph)
	win.Resize(fyne.NewSize(50, 70))
	defer win.Close()

	graph.Resize(fyne.NewSize(50, 70))
	graph.SetData(data)
	img := makeRasterize(win, graph)

	assertSize(t, img, graph)

}

func TestLineChart_RasterizerWithNegative(t *testing.T) {
	graph := createTestLineChart()
	data := []float64{-1, -2, -3, -4, -5, -6, -7, -8, -9, -10}
	graph.SetData(data)

	win := test.NewWindow(graph)
	win.Resize(fyne.NewSize(500, 300))
	defer win.Close()

	graph.Resize(fyne.NewSize(500, 300))

	img := makeRasterize(win, graph)
	assertSize(t, img, graph)

	data = []float64{-5, -4, -3, -2, -1, 0, 1, 2, 3, 4}
	graph.SetData(data)
	graph.Resize(fyne.NewSize(500, 300))
	img = makeRasterize(win, graph)
	assertSize(t, img, graph)
}

func TestLineChart_WithNoData(t *testing.T) {
	graph := createTestLineChart()
	win := test.NewWindow(graph)
	win.Resize(fyne.NewSize(500, 300))
	defer win.Close()

	assert.Len(t, graph.data, 0)
	assert.Equal(t, graph.opts.StrokeColor, theme.ForegroundColor())
	assert.Equal(t, graph.opts.FillColor, theme.DisabledButtonColor())
	assert.Equal(t, graph.opts.StrokeWidth, float32(1))

	// call rasterizer
	img := makeRasterize(win, graph)
	assertSize(t, img, graph)
}

func TestLineChart_GetOpts(t *testing.T) {
	opts := &PolygonCharthOpts{
		StrokeColor: color.RGBA{0x11, 0x22, 0x33, 255},
		FillColor:   color.RGBA{0x44, 0x55, 0x66, 255},
		StrokeWidth: 5,
	}
	graph := NewLineChart(opts)

	assert.Equal(t, graph.opts, opts)
	// in case of, check all fields
	assert.Equal(t, graph.opts.StrokeColor, color.RGBA{0x11, 0x22, 0x33, 255})
	assert.Equal(t, graph.opts.FillColor, color.RGBA{0x44, 0x55, 0x66, 255})
	assert.Equal(t, graph.opts.StrokeWidth, float32(5))
}

func TestLineChart_GetValAndCurvePos(t *testing.T) {
	graph := createTestLineChart()
	win := test.NewWindow(graph)
	win.Resize(fyne.NewSize(500, 300))
	defer win.Close()

	graph.CreateRenderer()
	graph.Resize(fyne.NewSize(500, 300))

	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	graph.SetData(data)
	graph.rasterize(500, 300)

	// Get the value at the center of the graph
	x, y := graph.GetDataPosAt(fyne.NewPos(289, 200))
	assert.Equal(t, float64(6), x)
	assert.Equal(t, float32(250), y.X)
	assert.Equal(t, float32(120), y.Y)
}
