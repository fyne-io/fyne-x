package main

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/widget/charts"
)

// SliderGraph is a LineChart with a Slider below to "scale" the data.
type SliderGraph struct {
	*charts.LineChart

	container *fyne.Container
	slider    *widget.Slider
	precision float64
}

// NewSliderGraph returns a SliderGraph.
func NewSliderGraph(nMin, nMax float64, precision float64) *SliderGraph {

	chart := charts.NewLineChart(nil)
	slider := widget.NewSlider(float64(nMin), float64(nMax))
	container := container.NewBorder(nil, slider, nil, nil, chart)

	// use border layout
	sg := &SliderGraph{
		container: container,
		LineChart: chart,
		slider:    slider,
		precision: precision,
	}

	//sg.LineChart.SetData(siny)
	sg.slider.OnChanged = sg.slided // connect slider event to "slided"

	// place the cursor to 20% (and so, init data)
	sg.slider.SetValue(nMin + (nMax-nMin)*.2)

	return sg
}

// Container returns the container of the widget.
func (sg *SliderGraph) Container() fyne.CanvasObject {
	return sg.container
}

// change / scale the dataset.
func (sg *SliderGraph) slided(value float64) {
	siny := make([]float64, int(value))
	for i := 0; i < int(value); i++ {
		siny[i] = math.Sin(float64(i) / float64(sg.precision))
	}
	sg.LineChart.SetData(siny)
	sg.LineChart.Refresh()
}
