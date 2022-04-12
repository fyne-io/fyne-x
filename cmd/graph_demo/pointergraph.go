package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/x/fyne/widget/charts"
)

// LineChartWithMouse is a line chart with mouse events handlers. This example
// draws a pointer + 2 lines over the chart when the mous is over it.
type LineChartWithMouse struct {
	*charts.LineChart
	mouseOver bool
}

// NewLineChartWithMouse creates a new line chart with mouse events handlers.
func NewLineChartWithMouse() *LineChartWithMouse {
	return &LineChartWithMouse{
		LineChart: charts.NewLineChart(nil),
		mouseOver: false,
	}
}

// IsMouseOver is true when the mouse is over the graph
func (c *LineChartWithMouse) IsMouseOver() bool {
	return c.mouseOver
}

// MouseDown is called when the mouse enters the chart. Here we only forward
// the "moved" event handler to draw the pointer + 2 lines over the chart.
func (c *LineChartWithMouse) MouseIn(e *desktop.MouseEvent) {
	c.mouseOver = true
	c.MouseMoved(e)
}

// MouseOut is called when the mouse is out of the chart. Here we remove the
// pointer + lines from the graph.
func (c *LineChartWithMouse) MouseOut() {
	c.mouseOver = false
	drawZone := c.GetDrawable()
	drawZone.Objects = nil
}

// MouseMoved is called when the mouse is moved over the chart.
// Here we draw the pointer as a circle under the mouse on the curve and 2 lines.
func (c *LineChartWithMouse) MouseMoved(e *desktop.MouseEvent) {

	// get the value and the curve position to the nearest found data
	val, curvePos := c.GetDataPosAt(e.Position)

	// prepare the vertical and horizontal lines to draw
	lineColor := theme.DisabledColor()
	verticalLine := canvas.NewLine(lineColor)
	verticalLine.Position1 = fyne.NewPos(curvePos.X, 0)
	verticalLine.Position2 = fyne.NewPos(curvePos.X, c.Size().Height)

	horizontalLine := canvas.NewLine(lineColor)
	horizontalLine.Position1 = fyne.NewPos(0, curvePos.Y)
	horizontalLine.Position2 = fyne.NewPos(c.Size().Width, curvePos.Y)

	// place a circle on the curve, we apply an offset on the circle corresponding to its size/2
	circle := canvas.NewCircle(theme.ForegroundColor())
	circle.Resize(fyne.NewSize(theme.TextSize()*.5, theme.TextSize()*.5))
	circle.Move(curvePos.Subtract(fyne.NewPos(
		circle.Size().Height/2, circle.Size().Width/2,
	)))

	// display the value over the circle
	text := canvas.NewText(fmt.Sprintf("%.02f", val), theme.ForegroundColor())
	text.TextSize = theme.TextSize() * 0.7
	text.Move(curvePos.Add(fyne.NewPos(
		text.TextSize*1.5, -text.TextSize*0.5,
	)))

	// then add line, circle and text to the graph
	drawZone := c.GetDrawable()
	drawZone.Objects = []fyne.CanvasObject{verticalLine, horizontalLine, circle, text}
	c.Refresh()
}
