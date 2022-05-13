package main

import (
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/widget/charts"
)

func main() {
	app := app.New()
	w := app.NewWindow("Graphs")

	var graphWidgets []fyne.CanvasObject

	// Create a basic LineChart
	lineChart := charts.NewLineChart(nil)
	lineChart.SetData([]float64{1, 3, -2, -4, 0, 4, 5, 6, 4})
	graphWidgets = append(graphWidgets, container.NewBorder(
		widget.NewLabelWithStyle("LineChart", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		lineChart,
	))

	// create n graphs
	g := NewLineChartWithMouse()
	animateLineChart(g)

	// Set a title for the graph, use nice Border layout
	graphWidgets = append(graphWidgets, container.NewBorder(
		widget.NewLabelWithStyle("Custom mouse event", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		g,
	))

	// create a basic bar chart
	barChart := charts.NewBarChart(nil)
	barChart.SetData([]float64{1, 3, 6, 4, 5, 6, 4})
	graphWidgets = append(graphWidgets, container.NewBorder(
		widget.NewLabelWithStyle("Bar Chart", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		barChart,
	))

	// Create a lineChart with custom color
	lineChartWithColor := charts.NewLineChart(nil)
	lineChartWithColor.SetData([]float64{1, 3, -2, -4, 0, 4, 5, 6, 4})
	lineChartWithColor.Options().StrokeColor = theme.PrimaryColor()
	lineChartWithColor.Options().FillColor = theme.ButtonColor()
	graphWidgets = append(graphWidgets, container.NewBorder(
		widget.NewLabelWithStyle("LineChart with custom color", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		lineChartWithColor,
	))

	// create a barchart
	bar := charts.NewBarChart(nil)
	animateBarChart(bar)
	graphWidgets = append(graphWidgets, container.NewBorder(
		widget.NewLabelWithStyle("Animated Bar Chart", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		bar,
	))

	// Create a overrided LineChart that has got a slider to scale the data
	sinContainer := NewSliderGraph(100, 1000, 10)
	graphWidgets = append(graphWidgets, container.NewBorder(
		widget.NewLabelWithStyle("Custom zommable", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		sinContainer.Container(),
	))

	// Greate the UI
	grid := container.NewGridWithColumns(2, graphWidgets...)
	w.SetContent(grid)
	w.Resize(fyne.NewSize(580, 340))
	w.ShowAndRun()
}

func animateLineChart(g *LineChartWithMouse) {
	g.Options().SetGraphRange(&charts.GraphRange{YMin: 0, YMax: 90})
	data := make([]float64, 64)
	for d := range data {
		data[d] = rand.Float64()*50 + 20
	}
	go func() {
		// Contiuously update the data

		// remove the first data point and add a new one each 500ms
		for range time.Tick(500 * time.Millisecond) {
			if g.IsMouseOver() {
				continue
			}
			data = append(data[1:], rand.Float64()*50+20)
			g.SetData(data)
			g.Refresh()
		}
	}()

}

func animateBarChart(chart *charts.BarChart) {
	chart.Options().SetGraphRange(&charts.GraphRange{YMin: 0, YMax: 90})
	initialColor := chart.Options().FillColor
	data := make([]float64, 10)
	for d := range data {
		data[d] = rand.Float64()*50 + 20
	}
	chart.SetData(data)
	chart.Refresh()

	go func() {
		// Contiuously update the data

		// remove the first data point and add a new one each 500ms
		for range time.Tick(500 * time.Millisecond) {
			data = append(data[1:], rand.Float64()*50+20)
			mean := 0.0
			for _, d := range data {
				mean += d
			}
			mean /= float64(len(data))

			// let's play with color, it changes when latest values are higher a certain threshold
			chart.Options().FillColor = initialColor
			if mean > 50 {
				chart.Options().FillColor = theme.ErrorColor()
			} else if mean > 45 {
				chart.Options().FillColor = theme.PrimaryColor()
			}
			chart.SetData(data)
			chart.Refresh()
		}
		// let's play with color
	}()

}
