package charts

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func createTestHistoChart() *BarChart {
	return NewBarChart(nil)
}

func createTestHistoChartWithOptions() *BarChart {
	return NewBarChart(&LineCharthOpts{
		StrokeColor: color.RGBA{0x11, 0x22, 0x33, 255},
		FillColor:   color.RGBA{0x44, 0x55, 0x66, 255},
		StrokeWidth: 5,
	})
}
func TestHistogram_Create(t *testing.T) {
	histo := createTestHistoChart()
	win := test.NewWindow(histo)
	win.Resize(fyne.NewSize(500, 300))
	defer win.Close()
}
