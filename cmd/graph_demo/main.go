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

	graphWidgets := make([]fyne.CanvasObject, 0)

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
	animateData(g)

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
	lineChartWithColor.GetOptions().StrokeColor = theme.PrimaryColor()
	lineChartWithColor.GetOptions().FillColor = theme.ButtonColor()
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

func animateData(g *LineChartWithMouse) {
	g.SetGraphRange(&charts.GraphRange{YMin: -1, YMax: 150})
	data := make([]float64, 64)
	for d := range data {
		data[d] = rand.Float64()*100 + 20
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
		}
	}()

}

func animateBarChart(chart *charts.BarChart) {
	chart.SetGraphRange(&charts.GraphRange{YMin: 0, YMax: 150})
	initialColor := chart.GetOptions().FillColor
	data := make([]float64, 10)
	for d := range data {
		data[d] = rand.Float64()*100 + 20
	}

	go func() {
		// Contiuously update the data

		// remove the first data point and add a new one each 500ms
		for range time.Tick(500 * time.Millisecond) {
			data = append(data[1:], rand.Float64()*50+20)
			max := 0.0
			for _, v := range data[3:] {
				if v > max {
					max = v
				}
			}
			// let's play with color, it changes when latest values are higher a certain threshold
			chart.GetOptions().FillColor = initialColor
			if max > 45 {
				chart.GetOptions().FillColor = theme.ErrorColor()
			} else if max > 35 {
				chart.GetOptions().FillColor = theme.PrimaryColor()
			}
			chart.SetData(data)
		}
		// let's play with color
	}()

}
